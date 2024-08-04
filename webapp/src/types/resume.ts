export type ResumeItem = {
    id: string;
    customerId: string;
    title: string;
    createdAt: Date;
    updatedAt: Date;
};

export type ResumeAbout = {
    resume_id: string;
    name: string;
    email: string;
    phone: string;
    title: string;
    location: string;
    github: string;
    linkedin: string;
    created_at: Date;
}

export type ResumeChecklistItem = {
    completed: boolean
    message: string
}

export type ResumeApplicationStatus =
    | 'not-started'
    | 'in-progress'
    | 'generated'
    | 'applied'
    | 'heard-back'
    | 'interviewing'
    | 'job-offer'
    | 'accepted';

export type ResumeApplication = {
    id: string;
    resumeId: string;
    title: string;
    link: string;
    companySite: string;
    rawText: string;
    status: ResumeApplicationStatus;
    createdAt: Date;
    updatedAt: Date;
}



export type CreateResumeApplicationRequest = {
    title: string,
    link: string,
    companySite: String,
    rawText: string,
}