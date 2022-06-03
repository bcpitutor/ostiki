# tikiserver

**Tiki** is a *limited life tickets* based authentication mechanism to
manage user authentication for the cloud operators. It employs the
client-server model and **Tikiserver** is the server component of it.

**Tikiserver** offers a REST API for clients to be able to connect it.

## What is a ticket?

A ticket in this context, is a data structure that contains one or
more sensitive information (usually login credentials or similar
secret information) made available to a user, backed with a 3rd party
authentication system (such as Google Workspaces) temporarily.

An example JSON representation of this data structure is similar to the below:

```
{
 "TicketPath": "organization/division/region/aws/s3-master",
 "TicketType": "awsTicket",
 "AwsAssumeRole": {
   "RoleArn": "arn:aws:iam::account_number:role/RoleName",
   "Ttl": 3600
 },
 "AwsPermissions": {
   "Action": [
      "s3:*",
   ],
   "Effect": "Allow",
   "Resource": "*"
 },
 "CreatedAt": "1642805843",
 "CreatedBy": "some@email.address",
 "OwnersGroup": [
    "s3-admins"
 ],
 "TicketInfo": "This ticket is for to S3 admins",
 "TicketRegion": "us-west-1",
 "UpdatedAt": "1642805843",
 "UpdatedBy": "some@email.address"
}

```

This example ticket, is a awsTicket (being used to manage AWS
Resources) and grants the permissions of the role specified in RoleArn
field at specified region, to the members of "s3-admins" group for one
hour. By using this ticket, members of the specified group can manage
s3 resources without logging onto any AWS accounts (root or IAM
users).

Similarily, different types of tickets can carry different types of
credentials or sensitive information.

Organizations can use Tiki to easily manage their users' cloud
resources permissions without creating personal users on the cloud
operator. All they need to create is their users Google Workspace
emails. Tiki can track individuals' ticket obtention even though they
don't have individual user accounts.
