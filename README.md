# go-and-databases
A test application for the blog post

In order to show the difference between using PGPool, Postgres and using Go's limits on databases I have created this small user service.

The most important part of this is that the `limits` can be turned on and off via the `-limits` flag and the `DSN` can be manipulated to connect either to PGPool or Postgres directly.

I have added run configs for Jetbrains based IDEs but the information in them is valuable for any way you want to run them. `docker-compose` will need to run first, then manually add the `users` table using the `sql` in the sql folder.
