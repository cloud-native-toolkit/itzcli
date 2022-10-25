package dr

const CascTemplateString = `
unclassified:
  location:
    url: http://server_ip:8080/
jenkins:
  securityRealm:
    local:
      allowsSignup: false
      users:
        - id: ${JENKINS_ADMIN_ID}
          password: ${JENKINS_ADMIN_PASSWORD}
  authorizationStrategy: unsecured
`
