pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.24.3'
    }
    
    stages {
        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }
        
        stage('Setup Go Environment') {
            steps {
                sh 'go version'
                sh 'export GO111MODULE=on'
            }
        }
        
        stage('Download Dependencies') {
            steps {
                sh 'go mod download'
                sh 'go mod verify'
            }
        }
        
        stage('Build') {
            steps {
                sh 'go build -v -o pakaidonk-backend .'
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