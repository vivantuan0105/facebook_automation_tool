export namespace models {
	
	export class DashboardStats {
	    total_posts: number;
	    total_tasks: number;
	    successful_runs: number;
	    recent_errors: number;
	    connection_status: string;
	    connected_account: string;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_posts = source["total_posts"];
	        this.total_tasks = source["total_tasks"];
	        this.successful_runs = source["successful_runs"];
	        this.recent_errors = source["recent_errors"];
	        this.connection_status = source["connection_status"];
	        this.connected_account = source["connected_account"];
	    }
	}
	export class Log {
	    id: number;
	    // Go type: time
	    time: any;
	    module: string;
	    action: string;
	    status: string;
	    details: string;
	
	    static createFrom(source: any = {}) {
	        return new Log(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.time = this.convertValues(source["time"], null);
	        this.module = source["module"];
	        this.action = source["action"];
	        this.status = source["status"];
	        this.details = source["details"];
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
	export class Post {
	    id: number;
	    title: string;
	    content: string;
	    summary: string;
	    status: string;
	    notes: string;
	    // Go type: time
	    post_at: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Post(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.summary = source["summary"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.post_at = this.convertValues(source["post_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Settings {
	    id: number;
	    theme: string;
	    app_name: string;
	    database_path: string;
	    auto_start: boolean;
	    log_level: string;
	    timezone: string;
	    connection_status: string;
	    connected_account: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.theme = source["theme"];
	        this.app_name = source["app_name"];
	        this.database_path = source["database_path"];
	        this.auto_start = source["auto_start"];
	        this.log_level = source["log_level"];
	        this.timezone = source["timezone"];
	        this.connection_status = source["connection_status"];
	        this.connected_account = source["connected_account"];
	    }
	}
	export class Task {
	    id: number;
	    name: string;
	    description: string;
	    type: string;
	    schedule: string;
	    enabled: boolean;
	    // Go type: time
	    last_run?: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.schedule = source["schedule"];
	        this.enabled = source["enabled"];
	        this.last_run = this.convertValues(source["last_run"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

