export namespace main {
	
	export class Timing {
	    Full: string;
	    Half: string;
	    Quarter: string;
	    Eighth: string;
	    Sixteenth: string;
	    ThirtySecond: string;
	    SixtyFourth: string;
	    OneTwentyEighth: string;
	
	    static createFrom(source: any = {}) {
	        return new Timing(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Full = source["Full"];
	        this.Half = source["Half"];
	        this.Quarter = source["Quarter"];
	        this.Eighth = source["Eighth"];
	        this.Sixteenth = source["Sixteenth"];
	        this.ThirtySecond = source["ThirtySecond"];
	        this.SixtyFourth = source["SixtyFourth"];
	        this.OneTwentyEighth = source["OneTwentyEighth"];
	    }
	}
	export class UpdateInfo {
	    available: boolean;
	    message: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.message = source["message"];
	        this.url = source["url"];
	    }
	}

}

