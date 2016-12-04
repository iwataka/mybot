# API Design

+ /api/config/

    _Get or change the configuration of this app_

    + GET /api/config/
        - request
            * body: nothing
        - response
            * body: JSON-formatted configuration data

    + POST /api/config/
        - request
            * body: JSON-formatted configuration data
        - response
            * body: nothing

+ /api/features/

    _Get or change the state of each feature of this app_

    + GET /api/features/twitter/listen/myself
        - request
            * body: nothing
        - response
            * body: result
            ```
            {
                status: true,
            }
            ```

    + POST /api/features/twitter/listen/myself
        - request
            * body: If you want to activate this feature, like this:
            ```
            {
                status: true,
            }
            ```
        - response
            * body: nothing

    + GET /api/features/twitter/listen/users
        - request
            * body: nothing
        - response
            * body: result
            ```
            {
                status: true,
            }
            ```

    + POST /api/features/twitter/listen/users
        - request
            * body: If you want to activate this feature, like this:
            ```
            {
                status: true,
            }
            ```
        - response
            * body: nothing

    + GET /api/features/twitter/periodic
        - request
            * body: nothing
        - response
            * body: result
            ```
            {
                status: true,
            }
            ```

    + POST /api/features/twitter/periodic
        - request
            * body: If you want to activate this feature, like this:
            ```
            {
                status: true,
            }
            ```
        - response
            * body: nothing

    + GET /api/features/github/periodic
        - request
            * body: nothing
        - response
            * body: result
            ```
            {
                status: true,
            }
            ```

    + POST /api/features/github/periodic
        - request
            * body: If you want to activate this feature, like this:
            ```
            {
                status: true,
            }
            ```
        - response
            * body: nothing

    + GET /api/features/others/monitor/config
        - request
            * body: nothing
        - response
            * body: result
            ```
            {
                status: true,
            }
            ```

    + POST /api/features/others/monitor/config
        - request
            * body: If you want to activate this feature, like this:
            ```
            {
                status: true,
            }
            ```
        - response
            * body: nothing
