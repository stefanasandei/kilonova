import { getCall } from "../net";

export type Submission = {
	id: number;
	created_at: string;
	max_time: number;
	max_memory: number;
	score: number;
	status: string;
	code_size: number;
};

export type UserBrief = {
	id: number;
	name: string;
	admin: boolean;
	proposer: boolean;
	bio: string;
};

export type Problem = {
	id: number;
	name: string;
};

export type ResultSubmission = {
	sub: Submission;
	author: UserBrief;
	problem: Problem;
	hidden: boolean;
};

export type SubmissionQuery = {
	user_id?: number;
	problem_id?: number;
	score?: number;
	status?: string;
	lang?: string; // TODO: allow undefined?

	page: number;

	compile_error?: boolean;
	ordering?: string;
	ascending?: boolean;

	limit?: number;
};

function serializeQuery(q: SubmissionQuery): object {
	return {
		ordering: typeof q.ordering !== "undefined" ? q.ordering : "id",
		ascending: (typeof q.ordering !== "undefined" && q.ascending) || false,
		user_id: typeof q.user_id !== "undefined" && q.user_id > 0 ? q.user_id : undefined,
		problem_id: typeof q.problem_id !== "undefined" && q.problem_id > 0 ? q.problem_id : undefined,
		status: q.status !== "" ? q.status : undefined,
		score: typeof q.score !== "undefined" && q.score >= 0 ? q.score : undefined,
		lang: typeof q.lang !== "undefined" && q.lang !== "" ? q.lang : undefined,
		compile_error: q.compile_error,
		offset: (q.page - 1) * 50,
		limit: typeof q.limit !== "undefined" && q.limit > 0 ? q.limit : 50,
	};
}

export async function getSubmissions(q: SubmissionQuery) {
	let res = await getCall<{
		count: number;
		subs: ResultSubmission[];
	}>("/submissions/get", serializeQuery(q));
	if (res.status === "error") {
		throw new Error(res.data);
	}
	return res.data;
}