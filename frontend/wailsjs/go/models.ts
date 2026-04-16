export namespace models {
	
	export class ActivityLog {
	    id: number;
	    time: string;
	    module: string;
	    action: string;
	    status: string;
	    details: string;
	
	    static createFrom(source: any = {}) {
	        return new ActivityLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.time = source["time"];
	        this.module = source["module"];
	        this.action = source["action"];
	        this.status = source["status"];
	        this.details = source["details"];
	    }
	}
	export class AppSettings {
	    appName: string;
	    theme: string;
	    defaultExecutionMode: string;
	    maskCookieByDefault: boolean;
	    persistToJson: boolean;
	    graphqlLoginDocId: string;
	    graphqlLikeDocId: string;
	    graphqlCommentDocId: string;
	    graphqlPostDocId: string;
	    graphqlScanPostDocId: string;
	    graphqlProfileDocId: string;
	    graphqlFriendsDocId: string;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.theme = source["theme"];
	        this.defaultExecutionMode = source["defaultExecutionMode"];
	        this.maskCookieByDefault = source["maskCookieByDefault"];
	        this.persistToJson = source["persistToJson"];
	        this.graphqlLoginDocId = source["graphqlLoginDocId"];
	        this.graphqlLikeDocId = source["graphqlLikeDocId"];
	        this.graphqlCommentDocId = source["graphqlCommentDocId"];
	        this.graphqlPostDocId = source["graphqlPostDocId"];
	        this.graphqlScanPostDocId = source["graphqlScanPostDocId"];
	        this.graphqlProfileDocId = source["graphqlProfileDocId"];
	        this.graphqlFriendsDocId = source["graphqlFriendsDocId"];
	    }
	}
	export class DashboardStats {
	    totalTasks: number;
	    pending: number;
	    running: number;
	    success: number;
	    failed: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalTasks = source["totalTasks"];
	        this.pending = source["pending"];
	        this.running = source["running"];
	        this.success = source["success"];
	        this.failed = source["failed"];
	    }
	}
	export class FacebookAccount {
	    uid: string;
	    name: string;
	    cookie: string;
	    status: string;
	    gender: string;
	    birthday: string;
	    birthYear: string;
	    location: string;
	    friends: string;
	    posts: string;
	    followers: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new FacebookAccount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.name = source["name"];
	        this.cookie = source["cookie"];
	        this.status = source["status"];
	        this.gender = source["gender"];
	        this.birthday = source["birthday"];
	        this.birthYear = source["birthYear"];
	        this.location = source["location"];
	        this.friends = source["friends"];
	        this.posts = source["posts"];
	        this.followers = source["followers"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class LoginResult {
	    uid: string;
	    cookieFull: string;
	    status: string;
	    rawResponse: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.cookieFull = source["cookieFull"];
	        this.status = source["status"];
	        this.rawResponse = source["rawResponse"];
	    }
	}
	export class ReactionTask {
	    id: number;
	    cookie: string;
	    cookieMasked: string;
	    taskType: string;
	    targetMode: string;
	    postId: string;
	    postUrl: string;
	    reactionType: string;
	    message: string;
	    photoPaths: string[];
	    feedbackId: string;
	    feedbackReactionId: string;
	    actorId: string;
	    feedbackSource: string;
	    sessionId: string;
	    clientMutationId: string;
	    notes: string;
	    scheduleAt: string;
	    quantityLimit: number;
	    status: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ReactionTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.cookie = source["cookie"];
	        this.cookieMasked = source["cookieMasked"];
	        this.taskType = source["taskType"];
	        this.targetMode = source["targetMode"];
	        this.postId = source["postId"];
	        this.postUrl = source["postUrl"];
	        this.reactionType = source["reactionType"];
	        this.message = source["message"];
	        this.photoPaths = source["photoPaths"];
	        this.feedbackId = source["feedbackId"];
	        this.feedbackReactionId = source["feedbackReactionId"];
	        this.actorId = source["actorId"];
	        this.feedbackSource = source["feedbackSource"];
	        this.sessionId = source["sessionId"];
	        this.clientMutationId = source["clientMutationId"];
	        this.notes = source["notes"];
	        this.scheduleAt = source["scheduleAt"];
	        this.quantityLimit = source["quantityLimit"];
	        this.status = source["status"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class TaskExecution {
	    id: number;
	    taskId: number;
	    mode: string;
	    status: string;
	    startedAt: string;
	    finishedAt: string;
	    errorMessage: string;
	    resultSummary: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskExecution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.taskId = source["taskId"];
	        this.mode = source["mode"];
	        this.status = source["status"];
	        this.startedAt = source["startedAt"];
	        this.finishedAt = source["finishedAt"];
	        this.errorMessage = source["errorMessage"];
	        this.resultSummary = source["resultSummary"];
	    }
	}

}

