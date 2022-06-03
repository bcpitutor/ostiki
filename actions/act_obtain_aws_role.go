package actions

// type AwsCredentials struct {
// 	AccessKeyId     string
// 	SecretAccessKey string
// 	SessionToken    string
// 	Region          string
// }

// func ObtainAWSRoleWithToken(token string, sessionName string, roleArn string, ttl int32, region string) (AwsCredentials, error) {
// 	var creds AwsCredentials

// 	cfg, err := config.LoadDefaultConfig(
// 		context.TODO(),
// 		config.WithRegion(region),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	client := sts.NewFromConfig(cfg)
// 	params := &sts.AssumeRoleWithWebIdentityInput{
// 		RoleArn:          aws.String(roleArn),
// 		RoleSessionName:  aws.String(sessionName),
// 		WebIdentityToken: aws.String(token),
// 		DurationSeconds:  aws.Int32(ttl),
// 	}

// 	resp, err := client.AssumeRoleWithWebIdentity(
// 		context.TODO(),
// 		params,
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return creds, err
// 	}

// 	creds.AccessKeyId = *resp.Credentials.AccessKeyId
// 	creds.SecretAccessKey = *resp.Credentials.SecretAccessKey
// 	creds.SessionToken = *resp.Credentials.SessionToken
// 	creds.Region = region

// 	return creds, nil
// }
