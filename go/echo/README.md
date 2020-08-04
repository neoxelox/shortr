# ✌️ This is the golang version of Shortr.
The actual code is inside **`main.go`** and my goal was to keep it as much as **`simple`** but also **`robust`** and **`fast`**. The other packages are just custom-made connectors for the database and other helpers.
It is recommended that you **`reuse the packages`** for other golang Shortr implementations.

It's built with the web framework [**`echo`**](https://echo.labstack.com/).
The connection to the database is made via the fastest postgres driver [**`pgx`**](https://github.com/jackc/pgx).
The default echo logger was a bit slow, so I implemented the fast [**`zerolog`**](https://github.com/rs/zerolog) (totally overkill).

You will have noticed that no test have been made. I may add them in the future, but for now they are overkill for this simple application.
However, these custom made packages will be updated and tested in my [**`microservice-template repository`**](https://github.com/Neoxelox/microservice-template) (which for this small project I did not follow it's patterns).
