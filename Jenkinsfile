pipeline {
    agent { label 'rhel7-8gb' }

    options {
        timestamps()
        timeout(time: 15, unit: 'MINUTES')
    }

    stages {
        stage("Build crw-osde2e image") {
            steps {
                sh """
                        make build-container
                    """

            }
        }
        stage("Push harness tests to docker repository") {
            steps {
                withCredentials([string(credentialsId: 'quay.io-crw-token', variable: 'QUAY_TOKEN')]) {
                    sh """
                        docker login -u="crw+crwci" --password ${QUAY_TOKEN} quay.io/crw
                        make push-container
                    """
                }
            }
        }
    }
}
