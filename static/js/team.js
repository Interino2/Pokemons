const drawButton = document.querySelector("#draw-team");
const teamResults = document.querySelector("#team-results");
const teamStatus = document.querySelector("#team-status");

drawButton.addEventListener("click", () => {
    const pokemons = Array.from(document.querySelectorAll(".pokemon-tile"));

    if (pokemons.length < 6) {
        teamStatus.textContent = "Necesitas al menos seis Pokémon en la colección para sortear un equipo.";
        teamResults.replaceChildren();
        return;
    }

    for (let index = pokemons.length - 1; index > 0; index -= 1) {
        const randomIndex = Math.floor(Math.random() * (index + 1));
        [pokemons[index], pokemons[randomIndex]] = [pokemons[randomIndex], pokemons[index]];
    }

    const team = document.createDocumentFragment();
    for (const pokemon of pokemons.slice(0, 6)) {
        const member = document.createElement("article");
        member.className = "team-member";

        const number = document.createElement("span");
        number.className = "team-number";
        number.textContent = `N.º ${pokemon.dataset.pokemonId.padStart(3, "0")}`;
        member.append(number);

        const imageContainer = document.createElement("span");
        imageContainer.className = "team-art";
        if (pokemon.dataset.pokemonImage) {
            const image = document.createElement("img");
            image.src = pokemon.dataset.pokemonImage;
            image.alt = pokemon.dataset.pokemonName;
            image.loading = "lazy";
            imageContainer.append(image);
        } else {
            const placeholder = document.createElement("span");
            placeholder.className = "image-placeholder";
            placeholder.setAttribute("aria-hidden", "true");
            placeholder.textContent = "?";
            imageContainer.append(placeholder);
        }
        member.append(imageContainer);

        const name = document.createElement("span");
        name.className = "team-name";
        name.textContent = pokemon.dataset.pokemonName;
        member.append(name);

        const type = document.createElement("span");
        type.className = "tile-type";
        type.textContent = pokemon.dataset.pokemonType;
        member.append(type);
        team.append(member);
    }

    teamResults.replaceChildren(team);
    teamStatus.textContent = "¡Equipo listo! Puedes volver a sortear cuando quieras.";
});
