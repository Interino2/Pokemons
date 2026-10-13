CREATE TABLE IF NOT EXISTS pokemons (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    tipo VARCHAR(50) NOT NULL,
    altura NUMERIC(8, 2) NOT NULL CHECK (altura > 0),
    peso NUMERIC(8, 2) NOT NULL CHECK (peso > 0),
    imagen TEXT NOT NULL DEFAULT '',
    descripcion VARCHAR(1000) NOT NULL DEFAULT '',
    pokedex_number INTEGER CHECK (pokedex_number BETWEEN 1 AND 151)
);
