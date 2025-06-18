export namespace games {
	
	export class GameMedia {
	    icon?: string;
	    hero?: string;
	    logo?: string;
	
	    static createFrom(source: any = {}) {
	        return new GameMedia(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.icon = source["icon"];
	        this.hero = source["hero"];
	        this.logo = source["logo"];
	    }
	}
	export class GameExternalID {
	    steam?: string;
	    grid_db?: string;
	
	    static createFrom(source: any = {}) {
	        return new GameExternalID(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.steam = source["steam"];
	        this.grid_db = source["grid_db"];
	    }
	}
	export class Game {
	    id: string;
	    slug?: string;
	    name: string;
	    description?: string;
	    external_ids?: GameExternalID;
	    media?: GameMedia;
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.slug = source["slug"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.external_ids = this.convertValues(source["external_ids"], GameExternalID);
	        this.media = this.convertValues(source["media"], GameMedia);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace scripting {
	
	export class LuaPlugin {
	    id: number[];
	    author: string;
	    name: string;
	    description: string;
	    version: string;
	    entry: string;
	    capabilities?: string[];
	    functionality?: string[];
	    settings: Record<string, plugins.SettingDefinition>;
	    PluginDir: string;
	    IsPacked: boolean;
	    Enabled: boolean;
	    Settings: Record<string, any>;
	    // Go type: utils
	    Cache?: any;
	
	    static createFrom(source: any = {}) {
	        return new LuaPlugin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.author = source["author"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.version = source["version"];
	        this.entry = source["entry"];
	        this.capabilities = source["capabilities"];
	        this.functionality = source["functionality"];
	        this.settings = this.convertValues(source["settings"], plugins.SettingDefinition, true);
	        this.PluginDir = source["PluginDir"];
	        this.IsPacked = source["IsPacked"];
	        this.Enabled = source["Enabled"];
	        this.Settings = source["Settings"];
	        this.Cache = this.convertValues(source["Cache"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

