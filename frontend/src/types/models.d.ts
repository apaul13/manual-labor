
export interface PaginationDetails {
    Limit: number;
    Offset: number;
}

export interface YearsResponse {
    Years: Year[];
    PaginationDetails: PaginationDetails;
}

export interface MakesResponse {
    Makes: Make[];
    PaginationDetails: PaginationDetails;
}

export interface ModelsResponse {
    Models: Model[];
    PaginationDetails: PaginationDetails;
}

export interface TrimsResponse {
    Trims: Trim[];
    PaginationDetails: PaginationDetails;
}

export interface Year {
    Year: string;
}

export interface Make {
    ID: number;
    Name: string;
    Year: string;
}

export interface Model {
    ID: number;
    Name: string;
    MakeId: number;
}

export interface Trim {
    ID: number;
    Name: string;
    ModelId: number;
}

export interface Car {
    MakeId: number;
    ModelId: number;
    TrimId: number;
}