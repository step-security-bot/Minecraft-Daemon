pipeline {
  agent any
  stages {
    stage('Setup') {
      steps {
        sh 'git submodule update --init --recursive'
        sh 'mkdir build'
        dir(path: 'build') {
          sh 'cmake -DCODE_COVERAGE=ON -DCMAKE_BUILD_TYPE=Debug ..'
        }

      }
    }

    stage('Build') {
      steps {
        dir(path: 'build') {
          sh 'make'
        }

      }
    }

    stage('Save Executable') {
      steps {
        archiveArtifacts(artifacts: 'build/src/Minecraft-Daemon', caseSensitive: true)
      }
    }

    stage('Test') {
      steps {
        dir(path: 'build') {
          sh 'ctest'
        }

      }
    }

    stage('Collect Coverage') {
      steps {
        sh './get_code_cov.sh'
        sh 'bash <(curl -s https://codecov.io/bash)'
      }
    }

  }
  environment {
    CODECOV_TOKEN = '5f5c6ae9-85f2-4a63-8219-69fef3151803'
  }
}