import { RouteHandlerMethod } from "fastify";
import axios from "axios";
import { User, AuthRouteHandlerMethod } from "./types";
import { PrismaClient } from "@prisma/client";

/**
 * Authenticates the user and executes the provided handler function.
 *
 * @param {AuthRouteHandlerMethod} handler - The handler function to execute after authentication.
 * @return {RouteHandlerMethod} The route handler method that handles the request and response.
 */
export const auth = (handler: AuthRouteHandlerMethod): RouteHandlerMethod => {
  return async (req, res) => {
    const token = req.headers.authorization?.split(" ")[1];
    const baseURL = process.env.USER_SERVICE;

    if (!token) {
      res.status(401);
      throw new Error("Not authenticated");
    }

    // Synchronous Messaging (Microservices Messaging Pattern)
    const result = await axios.get<User>(baseURL + "/myProfile", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    handler(req, res, result.data);
  };
};

/**
 * Checks if the current user is the owner of a playlist before executing a route handler method.
 *
 * @param {AuthRouteHandlerMethod} handler - The route handler method to be executed if the user is the owner.
 * @return {AuthRouteHandlerMethod} - The modified route handler method.
 */
export const isOwner = (
  handler: AuthRouteHandlerMethod
): AuthRouteHandlerMethod => {
  return async (req, res, user) => {
    const { id } = req.params as { id: string };
    const prisma = new PrismaClient();

    const playlist = await prisma.playlist.findFirst({
      where: { id },
    });

    if (playlist?.author !== user.ID) {
      throw new Error("You are not the owner of this playlist");
    }

    handler(req, res, user);
  };
};
