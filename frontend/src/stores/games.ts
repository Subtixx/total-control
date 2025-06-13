import {defineStore} from 'pinia';

export class Game {
    id: number;
    slug: string;
    name: string;
    description?: string;
    externalIds?: {
        grid_db?: string;
        steam?: string;
    }
    gamePath?: string;

    constructor(props: {
        id?: number;
        slug?: string;
        name?: string;
        description?: string;
        gamePath?: string;
        externalIds?: { grid_db?: string; steam?: string }
    }) {
        this.id = props.id || 0;
        this.slug = props.slug || '';
        this.name = props.name || '';
        this.description = props.description;
        this.externalIds = props.externalIds || {};
        this.gamePath = props.gamePath || '';
    }


    fetchImages() {
        //curl 'https://www.steamgriddb.com/api/public/game/10052' -H 'Referer: https://www.steamgriddb.com/game/10052/icons'
        /*
        {"success":true,"data":{"platforms":{"steam":{"id":"427520","gameId":10052,"metadata":{"store_asset_mtime":1611916591,"library_capsule":"en","library_logo":"en","library_hero":"en","steam_release_date":1597395480,"original_release_date":null,"logo_position":"{\"pinned_position\":\"CenterCenter\",\"width_pct\":\"53\",\"height_pct\":\"41\"}","header_image":"english","clienticon":"bcbf4ca8e6be3f0222e257396c8ba7dbc6caca29","icon":"267f5a89f36ab287e600a4e7d4e73d3d11f0fd7d","header_image_full":{"english":"header.jpg"},"library_capsule_full":{"image":{"english":"library_600x900.jpg"},"image2x":{"english":"library_600x900_2x.jpg"}},"library_hero_full":{"image":{"english":"library_hero.jpg"}},"library_logo_full":{"image":{"english":"logo.png"}}}},"gog":{"id":"1238653230","gameId":10052,"metadata":{"path":"\/game\/factorio","logo":null}}},"game":{"id":10052,"name":"Factorio","release_date":1456426816,"types":["steam","gog"],"verified":true},"header":{"external":false,"externa_url":null,"asset":{"id":11524,"style":"alternate","width":1920,"height":620,"nsfw":false,"humor":false,"notes":null,"language":"en","url":"https:\/\/cdn2.steamgriddb.com\/hero\/e2d52448d36918c575fa79d88647ba66.png","thumb":"https:\/\/cdn2.steamgriddb.com\/hero_thumb\/e2d52448d36918c575fa79d88647ba66.jpg","lock":false,"epilepsy":false,"upvotes":0,"downvotes":0,"downloads":110,"hearts":9,"can_vote":false,"date":1588630250,"mime":"image\/png","is_animated":false,"is_deleted":false,"animation_type":null,"processing":false,"show_boop":true,"author":{"name":"Morente","steam64":"76561197970305233","avatar":"https:\/\/avatars.steamstatic.com\/a4002879ba5e8dececed6b6f08596a6fe96df370.jpg","badges":[{"id":5,"name":"Traffic Conductor","description":"Contributed to the Steam Pre-Greenlight Project.","date":1602092823,"isHidden":false}]}}},"logo":{"external":false,"externa_url":null,"asset":{"id":1343,"style":"official","width":627,"height":105,"nsfw":false,"humor":false,"notes":null,"language":"en","url":"https:\/\/cdn2.steamgriddb.com\/logo\/674bfc5f6b72706fb769f5e93667bd23.png","thumb":"https:\/\/cdn2.steamgriddb.com\/logo_thumb\/674bfc5f6b72706fb769f5e93667bd23.png","lock":false,"epilepsy":false,"upvotes":0,"downvotes":0,"downloads":87,"hearts":4,"can_vote":false,"date":1576105797,"mime":"image\/png","is_animated":false,"is_deleted":false,"animation_type":null,"processing":false,"show_boop":true,"author":{"name":"stormyninja","steam64":"76561198067667175","avatar":"https:\/\/avatars.steamstatic.com\/d956bc4013e234d687e01e910297bd80cde68345.jpg","badges":[{"id":1,"name":"Patron","description":"Tier 1 Patreon supporter!","date":1748887466,"isHidden":false}]}}},"totals":{"grids":49,"heroes":11,"logos":13,"icons":2},"itad":null}}
        */
        const requestUrl = `https://www.steamgriddb.com/api/public/game/${this.externalIds?.grid_db}`;
        return fetch(requestUrl, {
            headers: {
                'Referer': `https://www.steamgriddb.com/game/${this.externalIds?.grid_db}/icons`
            }
        })
            .then(response => response.json())
            .then(data => {
                if (data.success && data.data) {
                    return data.data;
                } else {
                    throw new Error('Failed to fetch game images');
                }
            });
    }

    get icon() {
        return `/images/${this.slug}/icon.jpg`;
    }

    get capsule() {
        return `/images/${this.slug}/capsule.jpg`;
    }
}

const games: Game[] = [
    new Game({
        id: 1,
        slug: 'dying-light',
        name: 'Dying Light',
        description: 'An open world first-person survival horror game set in a post-apocalyptic world.',
        externalIds: {grid_db: '2716', steam: '239140'},
        gamePath: '/mnt/games/dying-light'
    }),
    new Game({
        id: 2,
        slug: 'dying-light-2',
        name: 'Dying Light 2',
        description: 'The sequel to Dying Light, featuring a larger world and more complex gameplay mechanics.',
        externalIds: {grid_db: '5148398', steam: '534380'},
        gamePath: '/mnt/games/dying-light-2'
    }),
    new Game({
        id: 3,
        slug: 'factorio',
        name: 'Factorio',
        description: 'A game about building and managing factories to produce items and automate processes.',
        externalIds: {grid_db: '10052', steam: '427520'},
        gamePath: '/mnt/games/factorio'
    }),
    new Game({
        id: 4,
        slug: 'satisfactory',
        name: 'Satisfactory',
        description: 'A first-person open-world factory building game with a focus on exploration and automation.',
        externalIds: {grid_db: '14065', steam: '526870'},
        gamePath: '/mnt/games/satisfactory'
    }),
    new Game({
        id: 5,
        slug: 'stardew-valley',
        name: 'Stardew Valley',
        description: 'A farming simulation game where players can grow crops, raise animals, and build relationships with villagers.',
        externalIds: {grid_db: '9569', steam: '413150'},
        gamePath: '/mnt/games/stardew-valley'
    }),
    new Game({
        id: 6,
        slug: 'terraria',
        name: 'Terraria',
        description: 'A 2D sandbox adventure game with crafting, building, and exploration elements.',
        externalIds: {grid_db: '1226', steam: '105600'},
        gamePath: '/mnt/games/terraria'
    }),
    new Game({
        id: 7,
        slug: 'starbound',
        name: 'Starbound',
        description: 'A procedurally generated space exploration game with crafting and building mechanics.',
        externalIds: {grid_db: '2048', steam: '211820'},
        gamePath: '/mnt/games/starbound'
    }),
];

export const useGamesStore = defineStore('games', {
    state: () => ({
        games: [] as Game[],
        recentGames: [] as Game[],
        loading: false,
        error: null as string | null,
    }),
    actions: {
        async addGame(game: Game) {
            this.loading = true;
            this.error = null;

            // For now, we simulate adding a game with a timeout
            try {
                await new Promise<void>(resolve => {
                    setTimeout(() => {
                        this.games.push(game);
                        resolve();
                    }, 500);
                });
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while adding the game.';
            }

            this.loading = false;
        },
        async fetchGames() {
            this.loading = true;
            this.error = null;

            // For now, we simulate fetching data with a timeout
            try {
                const response = await new Promise<Game[]>(resolve => {
                    setTimeout(() => {
                        resolve(games);
                    }, 1000);
                });

                this.games = response;
                this.recentGames = response.slice(0, 2); // Assume the first two are recent
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while fetching games.';
            }

            this.loading = false;
        },
        async fetchRecentGames() {
            this.loading = true;
            this.error = null;

            // For now, we simulate fetching data with a timeout
            try {
                const response = await new Promise<Game[]>(resolve => {
                    setTimeout(() => {
                        resolve(games.slice(0, 2)); // Assume the first two are recent
                    }, 1000);
                });

                this.recentGames = response;
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while fetching recent games.';
            }

            this.loading = false;
        },
        async fetchInstalledGames(): Promise<Game[]> {
            this.loading = true;
            this.error = null;

            // For now, we simulate fetching installed games with a timeout
            try {
                return await new Promise<Game[]>(resolve => {
                    setTimeout(() => {
                        resolve(games.filter(game => game.id % 2 === 0)); // Simulate installed games
                    }, 1000);
                });
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while fetching installed games.';
                return [];
            } finally {
                this.loading = false;
            }
        },
        async fetchGameById(id: number): Promise<Game | null> {
            this.loading = true;
            this.error = null;

            // For now, we simulate fetching a game by ID with a timeout
            try {
                const game = games.find(g => g.id === id) || null;
                if (game) {
                    return await new Promise<Game>(resolve => {
                        setTimeout(() => resolve(game), 500);
                    });
                } else {
                    throw new Error('Game not found');
                }
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while fetching the game.';
                return null;
            } finally {
                this.loading = false;
            }
        },
        async detectGame(path: string): Promise<Game | null> {
            this.loading = true;
            this.error = null;

            console.log(`Detecting game for path: ${path}`);

            // For now, we simulate detecting a game with a timeout
            try {
                const game = games.find(g => g.slug === path.split('/').pop()) || null;
                if (game) {
                    return await new Promise<Game>(resolve => {
                        setTimeout(() => resolve(game), 500);
                    });
                } else {
                    throw new Error('Game not found');
                }
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while detecting the game.';
                return null;
            } finally {
                this.loading = false;
            }
        },
        async searchGames(query: string): Promise<Game[]> {
            this.loading = true;
            this.error = null;

            // For now, we simulate searching games with a timeout
            try {
                const results = games.filter(game => game.name.toLowerCase().includes(query.toLowerCase()));
                return await new Promise<Game[]>(resolve => {
                    setTimeout(() => resolve(results), 2000);
                });
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'An error occurred while searching for games.';
                return [];
            } finally {
                this.loading = false;
            }
        }
    },
});
