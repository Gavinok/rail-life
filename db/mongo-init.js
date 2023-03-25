db = db.getSiblingDB("admin");
db.auth("admin", "admin");
db = db.getSiblingDB("social_media_app");

db.createUser({
  user: "user",
  pwd: "user",
  roles: [
    {
      role: "readWrite",
      db: "social_media_app",
    },
  ],
});

db.createCollection("users");
db.createCollection("posts");
db.createCollection("comments");
db.createCollection("notifications");
db.createCollection("activitys");
