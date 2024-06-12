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
                script {
                    def sourceCodePath = "${WORKSPACE}"
                    def sourceCodeChecksum = getSourceCodeChecksum(sourceCodePath)
                    echo "${sourceCodeChecksum} > $JENKINS_HOME/jobs/$JOB_NAME/$BUILD_NUMBER/sc_checksum"
                }
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

        stage('Build') {
            steps {
                sh 'make build'
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
                    // Loop through each app and generate checksums
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
