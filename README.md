
## Requirements

- Azure subscription
- Azure Active Directory B2C Tenant
    This content is ready to work with Azure Active Directory B2C, before start ensure you have your tenant ready.

    1. [Create B2C Tenant](https://docs.microsoft.com/en-us/azure/active-directory-b2c/tutorial-create-tenant)
    2. [Register an application and ROPC flow](https://docs.microsoft.com/en-us/azure/active-directory-b2c/add-ropc-policy?tabs=app-reg-ga&pivots=b2c-user-flow#ropc-flow-notes)

- Golang SDK 1.16.2 or higher.
- VS Code with Go Extension and Azure Functions extension.

## Solution

1. Function
    - Function service metadata (no code) used to initialize the Azure Function.
2. Handler
    - Handler service (golang service) used to add all the business logic to the Function.

Registered functions:

- authtoken - TimerTrigger to retrieve every 20 mins the authorization token used by graph internal API calls.
- authuser - HttpTrigger to authorize an AAD B2C user.
- createuser - HttpTrigger to create an AAD B2C user.
- getuser - HttpTrigger to obtain an AAD B2C user.
- publickeys - TimerTrigger to retrieve every day the public keys used to validate the jwt signature as part of the process validation.

## References

[Microsoft identity platform and OAuth 2.0 authorization code flow](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow)

# Initialize your development environment

- Configure handler service settings.

    All handler service settings are in the package: *identity\handler\settings*.

    Configure the following settings:

    ```
    // Tenant variable
    var Tenant = "[YOUR_TENANT_NAME].onmicrosoft.com"

    // TenantB2CLogin variable
    var TenantB2CLogin = "[YOUR_TENANT_NAME].b2clogin.com"

    // PrimaryTrustFrameworkPolicy variable
    var PrimaryTrustFrameworkPolicy = "[YOUR_ROPC_FLOW_NAME]"

    // ClientID variable
    var ClientID = "[YOUR_APP_ID]"

    // ClientSecret variable
    var ClientSecret = "[YOUR_APP_ID_SECRET]"
    ```

- Go to the handler folder and compile golang main module.

    ```
    handler> go build main.go
    ```

    *Note: if you are running in Windows you will get the "main.exe" file otherwise if you are on Linux or MacOS you will get the "main" file.

- Validate the host.json configuration in Azure Functions folder.

    If your development environment is Windows then let the defaultExecutablePath attribute as as "main.exe", otherwise, change it to: "main".

    ```
      "customHandler": {
            "description": {
            "defaultExecutablePath": "main.exe",
            "workingDirectory": "",
            "arguments": []
            },
            "enableForwardingHttpRequest": true
        }
    ```

- Move output handler file to Azure Functions folder.

    ```
    handler> mv main.exe ..\function\main.exe -Force
    ```

- Go to the function folder and run the Azure Function.

    ```
    function> func start
    ```

- If everything runs as expected you will obtain a terminal log like this.

    ```
    Azure Functions Core Tools
    Core Tools Version:       3.0.2996 Commit hash: c54cdc36323e9543ba11fb61dd107616e9022bba
    Function Runtime Version: 3.0.14916.0


    Functions:

            authuser: [POST] http://localhost:7071/api/authuser

            createuser: [POST] http://localhost:7071/api/createuser

            getuser: [GET] http://localhost:7071/api/getuser

            authtoken: timerTrigger

            publickeys: timerTrigger

    For detailed output, run func with --verbose flag.
    [2021-03-22T20:25:44.426Z] Executing 'Functions.authtoken' (Reason='Timer fired at 2021-03-22T14:25:44.4154475-06:00', Id=3ea7822e-c57c-4b68-9e3a-86e5f4b73919)
    [2021-03-22T20:25:44.427Z] Trigger Details: UnscheduledInvocationReason: IsPastDue, OriginalSchedule: 2021-03-22T14:20:00.0000000-06:00
    [2021-03-22T20:25:45.350Z] Worker process started and initialized.
    [2021-03-22T20:25:45.728Z] Outputs not set on http response for invocationId:3ea7822e-c57c-4b68-9e3a-86e5f4b73919
    [2021-03-22T20:25:45.728Z] ReturnValue not set on http response for invocationId:3ea7822e-c57c-4b68-9e3a-86e5f4b73919
    [2021-03-22T20:25:45.739Z] Executed 'Functions.authtoken' (Succeeded, Id=3ea7822e-c57c-4b68-9e3a-86e5f4b73919, Duration=1318ms)
    [2021-03-22T20:25:45.830Z] package: structs.common - initialized  
    [2021-03-22T20:25:45.831Z] package: structs.graph - initialized   
    [2021-03-22T20:25:45.831Z] package: structs.jwt - initialized     
    [2021-03-22T20:25:45.832Z] package: structs.requests - initialized
    [2021-03-22T20:25:45.832Z] package: settings - initialized        
    [2021-03-22T20:25:45.833Z] package: graph - initialized
    [2021-03-22T20:25:45.833Z] package: jwt - initialized
    [2021-03-22T20:25:45.834Z] package: handlers.http.authuser - initialized
    [2021-03-22T20:25:45.835Z] package: handlers.http.createuser - initialized
    [2021-03-22T20:25:45.835Z] package: handlers.http.getuser - initialized
    [2021-03-22T20:25:45.836Z] package: handlers.http.testgraph - initialized
    [2021-03-22T20:25:45.837Z] package: handlers.http.updateuser - initialized
    [2021-03-22T20:25:45.837Z] package: handlers.timer.authtoken - initialized
    [2021-03-22T20:25:45.838Z] Process of obtaining authorization token started
    [2021-03-22T20:25:45.839Z] Process of obtaining authorization token finalized
    [2021-03-22T20:25:45.840Z] package: handlers.timer.publickeys - initialized
    [2021-03-22T20:25:45.840Z] Process of obtaining public keys started
    [2021-03-22T20:25:45.841Z] Process of obtaining public keys finalized
    [2021-03-22T20:25:45.842Z] package: main - initialized
    [2021-03-22T20:25:45.842Z] 2021/03/22 14:25:45 About to listen on :2267. Go to https://127.0.0.1:2267/
    [2021-03-22T20:25:45.843Z] Process of obtaining authorization token started
    [2021-03-22T20:25:45.844Z] Process of obtaining authorization token finalized
    ```

    At this point all the golang packages has been initialized successfully by using the Azure Function process instance.

## Invoking APIs

- Function: createuser  
- Method: Post
- Url: http://localhost:7071/api/createuser
- Headers:
    - Content-Type: application/json
```
Sample payload:
{
    "givenName": "Roberto",
    "surname": "Cervantes",
    "email": "robece@microsoft.com",
    "password": "password"
}
```

- Function: authuser
- Method: Post
- Url: http://localhost:7071/api/authuser
- Headers:
    - Content-Type: application/json

```
Sample payload:
{
	"username": "robece@microsoft.com",
	"password": "password"
}
```

- Function: getuser
- Method: Get
- Url: http://localhost:7071/api/getuser?username=[username]
- Headers:
    - Authorization: Bearer [authToken]

```
Query parameter:
- username
```

## Azure Function deployment

There are two ways to deploy the function application, by: container or code, if container just build the source\identity\Dockerfile image, or by code just deploy Azure Function following the documentation [here](https://docs.microsoft.com/en-us/azure/azure-functions/create-first-function-vs-code-other?tabs=go%2Cwindows#publish-the-project-to-azure).