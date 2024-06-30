@Library('alvarium-pipelines') _

pipeline {
    agent any
    tools {
        go 'go-1.21'
    }
    environment {
        GO121MODULE = 'on'
        GIT_COMMIT = "c3212fc8e03d4e89a2971a279349726831872483"
    }
    stages {
        stage('prep - generate source code checksum') {
            steps {
                sh 'mkdir -p $JENKINS_HOME/jobs/$JOB_NAME/$BUILD_NUMBER/'
                // $PWD is the workspace dir (the cloned repo), this will generate 
                // an md5sum (checksum) for the repo and write it to `sc_checksum` in
                // the dir created above
                sh ''' find . -type f -exec md5sum {} + | LC_ALL=C sort | md5sum |\
                        cut -d" " -f1 \
                        > $JENKINS_HOME/jobs/$JOB_NAME/$BUILD_NUMBER/sc_checksum
                '''
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
                sh 'make cmd/transitor/transitor-demo'
            }
        }

        // stage('Dockerize') {
        //     steps {
        //         script {
        //             // Define the docker image names
        //             def appNames = ['creator', 'mutator', 'transitor']
        //             // Loop through each app and build the Docker image
        //             appNames.each { appName ->
        //                 def dockerImage = "${appName}-demo"
        //                 sh "docker build --build-arg TAG=${TAG} -t ${dockerImage} -f Dockerfile.${appName} ."
        //                 // TODO: push image to a registry
        //             }
        //         }
        //     }
        // }

        stage('alvarium - post-build annotations') {
            steps {
                script {
                    // Loop through each app and generate checksums
                    def appNames = ['transitor']
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
