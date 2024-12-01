CREATE TABLE songs (
    id serial primary key,
    name varchar(255) not null,
	artist varchar(255) not null,
	release_date date,
	text varchar(1000),
	link varchar(255)
)