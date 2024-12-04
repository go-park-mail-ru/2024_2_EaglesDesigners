db = db.getSiblingDB('admin');
db.auth('root', 'root');

db = db.getSiblingDB('files');
db.createUser({
    user: 'user',
    pwd: 'user',
    roles: [
        {
            role: 'readWrite',
            db: 'files',
        },
    ],
});

db.createCollection('mongodb_docker');