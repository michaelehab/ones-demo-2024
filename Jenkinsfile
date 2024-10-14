@Library('alvarium-pipelines') _

pipeline {
    agent any
    tools {
        go 'go-1.21'
    }
    environment {
        GO121MODULE = 'on'
        TAG = "${GIT_COMMIT}"
    }
    stages {
        stage('prep - generate source code checksum') {
            steps {
                sh 'mkdir -p $JENKINS_HOME/jobs/$JOB_NAME/$BUILD_NUMBER/'
                sh '''find . -type f -exec md5sum {} + |\
                        md5sum |\
                        cut -d" " -f1 \
                        > $JENKINS_HOME/jobs/$JOB_NAME/$BUILD_NUMBER/sc_checksum
                '''
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o cmd/creator/creator-demo ./cmd/creator'
                sh 'go build -o cmd/mutator/mutator-demo ./cmd/mutator'
                sh 'go build -o cmd/transitor/transitor-demo ./cmd/transitor'
            }
        }

        stage('alvarium - pre-build annotations') {
            steps {
                script {
                    def optionalParams = ['sourceCodeChecksumPath':"${JENKINS_HOME}/jobs/${JOB_NAME}/${BUILD_NUMBER}/sc_checksum"]
                    alvariumCreate(['source-code', 'vulnerability'], optionalParams)
                }
            }
        }

        stage('Dockerize') {
            steps {
                script {
                    // Define the docker image names
                    def appNames = ['creator', 'mutator', 'transitor']
                    // Loop through each app and build the Docker image
                    appNames.each { appName ->
                        def dockerImage = "${appName}-demo"
                        sh "docker build --build-arg TAG=${TAG} -t ${dockerImage} -f Dockerfile.${appName} ."
                        // TODO: push image to a registry
                    }
                }
            }
        }

        stage('alvarium - post-build annotations') {
            steps {
                script {
                    // Check if artifact has a valid checksum... Ideally this
                    // should be done by whatever is pulling the artifact but
                    // alvarium currently has no way of persisting information
                    // relating to annotator logic, which is why the checksum is
                    // being fetched from the file system instead of a persistent
                    // store
                    def appNames = ['creator', 'mutator', 'transitor']
                    appNames.each { appName ->
                        def artifactPath = "cmd/${appName}/${appName}-demo"
                        def checksumPath = "${JENKINS_HOME}/jobs/${JOB_NAME}/${BUILD_NUMBER}/${appName}-demo.checksum"
                        sh "md5sum ${artifactPath} | cut -d ' ' -f 1 | tr 'a-z' 'A-Z' | tr -d '\n' > ${checksumPath}"
                        
                        def optionalParams = [
                            "artifactPath": artifactPath,
                            "checksumPath": checksumPath
                        ]
                        alvariumTransit(['checksum'], optionalParams)
                    }
                }
            }
        }
    }
}
