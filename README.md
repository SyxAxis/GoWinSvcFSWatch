# GoWinSvcFSWatch
Pure Go coded Windows service that watches folders and alerts

Yet another Golang experiment. Using the core Golang libraries and the FSNotify project, could I get a Windows service to stand up and alert into Windows event log when files appeared in one or more folders. 

It's just an experiment but the options it opens up are interesting, use the service to notify other events or calls over REST, that a file has appeared, changed or been updated, eg file changes on a UNC path in network, post to KV store like DynamoDB or Couch, then send the file an S3 bucket, which kicks a Lambda off to process it, or just post the rawe file data to a REST server.


