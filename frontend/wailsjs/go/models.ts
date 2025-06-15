export namespace games {
	
	export class Game {
	    id: string;
	    slug?: string;
	    name: string;
	    steam_appid: number;
	    header_image?: string;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.slug = source["slug"];
	        this.name = source["name"];
	        this.steam_appid = source["steam_appid"];
	        this.header_image = source["header_image"];
	    }
	}

}

