pipeline {
    agent any
    
    stages {
        stage('Build') {
            steps {
                // Checkout code
                checkout scm
                
                // Check if Go is installed
                sh '''
                    echo "Checking for Go installation..."
                    if ! command -v go &> /dev/null; then
                        echo "ERROR: Go is not installed or not in PATH"
                        echo "Please install Go on the Jenkins agent"
                        exit 1
                    fi
                    
                    echo "Go version:"
                    go version
                    
                    echo "Setting up Go environment..."
                    export GO111MODULE=on
                    
                    echo "Downloading dependencies..."
                    go mod download
                    
                    echo "Building application..."
                    go build -v -o pakaidonk-backend .
                    
                    echo "Build completed successfully!"
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
            echo 'Please ensure Go is installed on the Jenkins agent'
        }
    }
}