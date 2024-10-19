INSERT INTO users (id, username, display_name, password_hash)
VALUES (gen_random_uuid(), 'admin',
        'Admin',
        '$2a$12$1ykeMAtpCbZD5tnp1hklIuXlo4KD4wD9Hu50.ph4BH7JqKJCKMGnC');