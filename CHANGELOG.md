# Changelog

## [Unreleased](https://github.com/monsoon/arc/tree/HEAD)

[Full Changelog](https://github.com/monsoon/arc/compare/2018.7/8...HEAD)

**Fixed bugs:**

- **Automation powershell script fails [\#145](https://gitHub.***REMOVED***/monsoon/arc/issues/145)**   
Trying to run a powersehll script on a windows machine fails with `%1 is not a valid Win32 application`.

## [2018.7/8](https://github.com/monsoon/arc/tree/2018.7/8) (2018-08-31)

[Full Changelog](https://github.com/monsoon/arc/compare/2018.5/6...2018.7/8)

**Implemented enhancements:**

- **Add cert expiration fact [\#140](https://gitHub.***REMOVED***/monsoon/arc/issues/140)**   
Add a new fact that reports the left hours to the cert expiration date.
- **Extend agent install script with the renew cert URL [\#138](https://gitHub.***REMOVED***/monsoon/arc/issues/138)**   
Extending the agent script with the renew cert URL will be needed / useful when initializing the arc node and afterwards when using the renewCert command or the auto renew cert function.
- **Check mosquitto and arc-api images for vulnerabilities [\#135](https://gitHub.***REMOVED***/monsoon/arc/issues/135)**   
Add a task in steps mosquitto-build and api-build to check for vulnerabilities in the pipeline.
- **Command to renew the node certs [\#131](https://gitHub.***REMOVED***/monsoon/arc/issues/131)**   
Write a command that renews the certs on the node side.
- **Optional Start API Server with TLS [\#127](https://gitHub.***REMOVED***/monsoon/arc/issues/127)**   
To be able to auto renew the certificates by the nodes we need first to run the API server with TLS. Running the Server with TLS we are able to extract the CN, OU and O from the certificate used on the tls communication.
- **Run integration tests against QA [\#125](https://gitHub.***REMOVED***/monsoon/arc/issues/125)**   
Decided to run the integration Tests against QA. Staging is not anymore reliable.
- **Upgrade go to v10 and change dependencies manager to dep [\#123](https://gitHub.***REMOVED***/monsoon/arc/issues/123)**   
- Upgrade to go v10 and compile the project to check everything works fine.
- **Initial Trust: One-Time Token & Launch-Index [\#49](https://gitHub.***REMOVED***/monsoon/arc/issues/49)**   
As far as I understand, the motivation for extending the openstack metadata server is, that you want to pass the same user-data into the batch-creation call and passing a \*\*one\*\*-time token in user-data for multiple machines does obviously not work.

## [2018.5/6](https://github.com/monsoon/arc/tree/2018.5/6) (2017-06-26)

[Full Changelog](https://github.com/monsoon/arc/compare/2018.3/4...2018.5/6)

## [2018.3/4](https://github.com/monsoon/arc/tree/2018.3/4) (2017-04-28)

[Full Changelog](https://github.com/monsoon/arc/compare/2018.1/2...2018.3/4)

## [2018.1/2](https://github.com/monsoon/arc/tree/2018.1/2) (2017-02-24)

[Full Changelog](https://github.com/monsoon/arc/compare/240d117ab77d707b3152473a30c19c604764ab00...2018.1/2)

**Closed issues:**

- **Agent state does not change to offline when terminating an instance. [\#96](https://gitHub.***REMOVED***/monsoon/arc/issues/96)**   
Check the heartbeat settings of the mqtt client



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*