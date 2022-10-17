export namespace main {
	
	export class returnStruct {
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new returnStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}

}

