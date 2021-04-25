db.createUser(
    {
      user: "admin",
      pwd: "password",
      roles: [ { role: "userAdmin", db: "realtime_chat" }, "readWrite" ]
    }
);