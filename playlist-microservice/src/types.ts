import { FastifyReply, FastifyRequest } from "fastify";

interface User {
  ID: string;
  FullName: string;
  Gender: string;
  Email: string;
  HashedPassword: string;
  Followers: number;
  Followings: number;
  ProfileURL: string;
  IsPrivate: boolean;
}

interface Music {
  ID: string;
  ArtistID: string;
  Likes: number;
  FileUrl: string;
  PosterUrl: string;
  Title: string;
  ShortDesc: string;
}

type AuthRouteHandlerMethod = (
  req: FastifyRequest,
  res: FastifyReply,
  user: User
) => {};

export { AuthRouteHandlerMethod, User, Music };
