# sqlboiler-bug

trying to repro it

- `./sqlboiler.sh db` stands up a configured db
- `./sqlboiler.sh test` runs a test against it (test is dumb, db must be fresh)
- `./sqlboiler.sh stop` stops the db
- `./sqlboiler.sh gen` regenerates the models according to the schema in `sample.sql`
