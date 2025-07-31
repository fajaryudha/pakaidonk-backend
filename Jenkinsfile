// Jenkinsfile (Ini harus ada di root repositori GitHub Anda)
pipeline {
    // --- PENTING: Gunakan agent docker ini di bagian paling atas pipeline ---
    agent {
        docker {
            image 'golang:1.22-alpine' // Image ini sudah berisi Go
            // Penting: Pasang docker.sock agar Docker di dalam container bisa berkomunikasi dengan Docker daemon host
            // Ini diperlukan untuk stage 'Build Docker Image' dan 'Deploy'
            args '-u 0 -v /var/run/docker.sock:/var/run/docker.sock' 
            // Tambahkan hak akses untuk Jenkins user di dalam container jika ada masalah permission
            // user 'root' // opsional, jika perlu root access di dalam container
        }
    }
    // ----------------------------------------------------------------------
    
    environment {
        APP_NAME = 'pakaidonk-backend'
        DOCKER_HUB_USERNAME = 'fajaryudha' // Ganti dengan username Docker Hub Anda
        DOCKER_HUB_REPO_NAME = 'pakaidonk-backend'
        DOCKER_HUB_CREDENTIALS_ID = 'fcbca270-544e-42e7-a52b-3e61253b2a10' // Ganti dengan ID kredensial Docker Hub Anda di Jenkins
        
        SSH_HOST = '8.215.8.98' // Ganti dengan host SSH Anda
        SSH_USERNAME = 'admin' // Ganti dengan username SSH Anda
        SSH_CREDENTIALS_ID = 'SSH' // Ganti dengan ID kredensial SSH Anda di Jenkins
        
        APP_VERSION = "1.0" // Menggunakan BUILD_NUMBER sebagai versi aplikasi
        // Pastikan Anda sudah membuat credentials di Jenkins (tipe Username & password untuk Docker Hub, SSH Username with private key untuk SSH)
    }

    stages {
        // Hapus stage 'Checkout SCM' yang redundant jika Anda punya, karena Jenkins sudah otomatis checkout.
        // Stage 'Declarative: Checkout SCM' akan tetap muncul otomatis di log oleh Jenkins.

        stage('Setup Go Environment') {
            steps {
                script {
                    echo "Verifying Go installation..."
                    // Baris ini akan berhasil jika agen docker berjalan dengan benar
                    sh 'go version' 
                    sh 'export GO111MODULE=on' // Pastikan Go Modules diaktifkan
                    echo "Go Environment is set up."
                }
            }
        }
        
        stage('Download Dependencies') {
            steps {
                script {
                    echo "Downloading Go dependencies..."
                    sh 'go mod download'
                    sh 'go mod verify'
                    echo "Dependencies downloaded."
                }
            }
        }
        
        stage('Build Go Application') {
            steps {
                script {
                    echo "Building Go application..."
                    // Gunakan CGO_ENABLED=0 untuk membuat binary statis (tanpa dependensi C)
                    // -ldflags "-s -w" mengurangi ukuran binary dengan menghapus info debug
                    sh 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ${APP_NAME} .'
                    echo "Go application built: ${APP_NAME}"
                    archiveArtifacts artifacts: "${APP_NAME}", fingerprinter: true // Simpan binary sebagai artefak Jenkins
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    echo "Building Docker image..."
                    // Instal docker-cli di dalam container Alpine agar perintah 'docker' dapat dijalankan
                    sh 'apk add --no-cache docker-cli' 
                    
                    sh "docker build -t ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:${APP_VERSION} ."
                    sh "docker tag ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:${APP_VERSION} ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:latest"
                    echo "Docker image built: ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:${APP_VERSION}"
                }
            }
        }

        stage('Push Docker Image') {
            when {
                expression { currentBuild.currentResult == 'SUCCESS' } // Hanya jalan jika tahap sebelumnya sukses
            }
            steps {
                script {
                    echo "Logging in to Docker Hub and pushing image..."
                    // Menggunakan kredensial dari Jenkins
                    withCredentials([usernamePassword(credentialsId: "${DOCKER_HUB_CREDENTIALS_ID}", passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        sh "echo \"${DOCKER_PASSWORD}\" | docker login -u \"${DOCKER_USERNAME}\" --password-stdin"
                        sh "docker push ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:${APP_VERSION}"
                        sh "docker push ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:latest"
                        sh "docker logout"
                    }
                    echo "Docker image pushed successfully."
                }
            }
        }

        stage('Deploy') {
            when {
                expression { currentBuild.currentResult == 'SUCCESS' } // Hanya jalan jika tahap sebelumnya sukses
            }
            steps {
                script {
                    echo "Deploying to ${SSH_HOST}..."
                    // Instal openssh-client di dalam container Alpine agar perintah 'ssh' tersedia
                    sh 'apk add --no-cache openssh-client' 
                    
                    // Menggunakan kredensial SSH dari Jenkins
                    withCredentials([sshUserPrivateKey(credentialsId: "${SSH_CREDENTIALS_ID}", keyFileVariable: 'SSH_KEY_FILE')]) {
                        // Perintah SSH untuk deployment
                        sh """
                            ssh -o StrictHostKeyChecking=no -i ${SSH_KEY_FILE} ${SSH_USERNAME}@${SSH_HOST} "
                                # Pastikan docker login di server tujuan jika perlu pull image private
                                echo \\"${DOCKER_PASSWORD}\\" | docker login -u \\"${DOCKER_USERNAME}\\" --password-stdin
                                
                                # Hentikan dan hapus container lama jika ada
                                if docker ps -a --format '{{.Names}}' | grep -q ${APP_NAME}; then
                                    echo 'Stopping and removing old container...'
                                    docker stop ${APP_NAME}
                                    docker rm ${APP_NAME}
                                fi
                                
                                # Pull image terbaru dari Docker Hub
                                echo 'Pulling latest image...'
                                docker pull ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:latest
                                
                                # Jalankan container baru
                                echo 'Starting new container...'
                                # Sesuaikan port mapping Anda (misal -p 80:8181)
                                docker run -d --name ${APP_NAME} -p 80:8181 ${DOCKER_HUB_USERNAME}/${DOCKER_HUB_REPO_NAME}:latest
                                
                                # Opsional: logout dari docker
                                docker logout
                                
                                echo 'Deployment complete.'
                            "
                        """
                    }
                    echo "Deployment to ${SSH_HOST} finished."
                }
            }
        }
    }
    
    post {
        always {
            cleanWs() // Membersihkan workspace Jenkins setelah build selesai
            echo "Pipeline finished. Status: ${currentBuild.currentResult}"
        }
        success {
            echo 'Build successful! 🎉'
        }
        failure {
            echo 'Build failed! 😭'
        }
        unstable {
            echo 'Build unstable, check warnings! 🚧'
        }
        aborted {
            echo 'Build aborted! 🛑'
        }
    }
}