pipeline {
    agent {
        docker {
            image 'golang:1.24.3-alpine'
        }
    }
    
    stages {
        stage('Build') {
            steps {
                sh '''
                    # Install git (required for go mod download)
                    apk add --no-cache git
                    
                    # Build the Go application
                    go mod download
                    go build -v -o pakaidonk-backend .
                '''
            }
        }
    }
    
    post {
        success {
            echo 'Build successful!'
        }
        failure {
            echo 'Build failed!'
        }
    }
}