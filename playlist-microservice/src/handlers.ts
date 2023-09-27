import { Playlist, PrismaClient } from "@prisma/client";
import { RouteHandlerMethod } from "fastify";
import { v4 as uuid, validate as validUUID } from "uuid";
import { AuthRouteHandlerMethod, Music } from "./types";
import { v2 as cloudinary } from "cloudinary";
import axios from "axios";

const prisma = new PrismaClient();

/**
 * Fetches the musics for a given playlist from the music service.
 *
 * @param {Playlist} playlist - The playlist object containing the list of music IDs.
 * @return {Promise<Playlist>} - The updated playlist object with the fetched musics.
 */
const fetchPlaylistMusics = async (playlist: Playlist) => {
  const musicServiceUrl = process.env.MUSIC_SERVICE!;

  const response = await axios.post<Music[]>(
    musicServiceUrl + "/retrieveMusicsByIds",
    {
      musicsIds: playlist.musics,
    }
  );

  return { ...playlist, musics: response.data };
};

/**
 * Retrieves a playlist by its ID and sends the result as a response.
 *
 * @param {Request} req - The request object.
 * @param {Response} res - The response object.
 * @return {Promise<void>} - A promise that resolves once the response is sent.
 */
const getPlaylistById: RouteHandlerMethod = async (req, res) => {
  const { id } = req.params as { id: string };

  if (id.length === 0 || !validUUID(id)) {
    throw new Error("Invalid ID");
  }

  const playlist = await prisma.playlist.findFirst({
    where: {
      id: id,
    },
  });

  if (!playlist) {
    throw new Error("Playlist not found");
  }

  const result = await fetchPlaylistMusics(playlist);

  res.status(200).send(result);
};

/**
 * Retrieves all playlists and their associated musics.
 *
 * @param {Request} req - The request object.
 * @param {Response} res - The response object.
 * @return {Promise<void>} The promise that resolves when the function completes.
 */
const getPlaylists: RouteHandlerMethod = async (req, res) => {
  const playlists = await prisma.playlist.findMany();

  if (!playlists) {
    throw new Error("Playlists not found or empty");
  }

  const result = await Promise.all(playlists.map(fetchPlaylistMusics));

  res.status(200).send(result);
};

/**
 * Like a playlist.
 *
 * @param {object} req - The request object.
 * @param {object} res - The response object.
 * @param {object} user - The user object.
 * @throws {Error} If the ID is invalid.
 * @throws {Error} If the playlist is not found.
 * @throws {Error} If the user has already liked the playlist.
 * @return {Promise<void>} A promise that resolves with no value.
 */
const likePlaylist: AuthRouteHandlerMethod = async (req, res, user) => {
  const { id } = req.params as { id: string };

  if (id.length === 0 || !validUUID(id)) {
    throw new Error("Invalid ID");
  }

  const playlist = await prisma.playlist.findFirst({
    where: {
      id,
    },
  });

  if (!playlist) {
    throw new Error("Playlist not found");
  }

  if (playlist.likes.includes(user.ID)) {
    throw new Error("You already liked this playlist");
  }

  await prisma.playlist.update({
    data: { likes: { push: user.ID } },
    where: { id },
  });

  res.status(200).send({ success: true });
};

/**
 * Unlike a playlist.
 *
 * @param {AuthRouteHandlerMethod} req - The request object.
 * @param {Response} res - The response object.
 * @param {User} user - The user object.
 * @throws {Error} If the ID is invalid.
 * @throws {Error} If the playlist is not found.
 * @throws {Error} If the user didn't like the playlist.
 * @return {Promise<void>} A promise that resolves when the playlist is unliked.
 */
const unlikePlaylist: AuthRouteHandlerMethod = async (req, res, user) => {
  const { id } = req.params as { id: string };

  if (id.length === 0 || !validUUID(id)) {
    throw new Error("Invalid ID");
  }

  const playlist = await prisma.playlist.findFirst({
    where: {
      id,
    },
  });

  if (!playlist) {
    throw new Error("Playlist not found");
  }

  if (!playlist.likes.includes(user.ID)) {
    throw new Error("You didn't like this playlist");
  }

  await prisma.playlist.update({
    data: { likes: playlist.likes.filter((id) => id !== user.ID) },
    where: { id },
  });

  res.status(200).send({ success: true });
};

/**
 * Uploads music to a playlist.
 *
 * @param {type} req - the request object
 * @param {type} res - the response object
 * @return {type} Promise<void> - a promise that resolves when the music is uploaded successfully
 */
const uploadMusicToPlaylist: AuthRouteHandlerMethod = async (req, res) => {
  const { musicId } = req.query as { musicId: string };
  const { id } = req.params as { id: string };

  if (musicId.length === 0 || !validUUID(musicId)) {
    throw new Error("Invalid ID");
  }

  if (id.length === 0 || !validUUID(id)) {
    throw new Error("Invalid ID");
  }

  const playlist = await prisma.playlist.findFirst({
    where: {
      id,
    },
  });

  if (!playlist) {
    throw new Error("Playlist not found");
  }

  await prisma.playlist.update({
    data: { musics: { push: musicId } },
    where: { id },
  });

  res.status(200).send({ success: true });
};

/**
 * Removes a music from a playlist.
 *
 * @param {string} req.query.musicId - The ID of the music to remove.
 * @param {string} req.params.id - The ID of the playlist.
 * @throws {Error} If the musicId is invalid or empty.
 * @throws {Error} If the id is invalid or empty.
 * @throws {Error} If the playlist is not found.
 * @return {void}
 */
const removeMusicFromPlaylist: AuthRouteHandlerMethod = async (req, res) => {
  const { musicId } = req.query as { musicId: string };
  const { id } = req.params as { id: string };

  if (musicId.length === 0 || !validUUID(musicId)) {
    throw new Error("Invalid ID");
  }

  if (id.length === 0 || !validUUID(id)) {
    throw new Error("Invalid ID");
  }

  const playlist = await prisma.playlist.findFirst({
    where: {
      id,
    },
  });

  if (!playlist) {
    throw new Error("Playlist not found");
  }

  await prisma.playlist.update({
    data: { musics: playlist.musics.filter((id) => id !== musicId) },
    where: { id },
  });

  res.status(200).send({ success: true });
};

/**
 * Creates a new playlist.
 *
 * @param {Object} req - The HTTP request object.
 * @param {Object} res - The HTTP response object.
 * @param {Object} user - The user object.
 * @return {Promise<void>} - A promise that resolves with no value.
 */
const createPlaylist: AuthRouteHandlerMethod = async (req, res, user) => {
  const { name, shortDesc, poster } = req.body as { [key: string]: string };

  const uploadedPoster = await cloudinary.uploader.upload(poster);

  const result = await prisma.playlist.create({
    data: {
      id: uuid(),
      name,
      shortDesc,
      author: user.ID,
      likes: [],
      musics: [],
      posterUrl: uploadedPoster.secure_url,
    },
    select: {
      id: true,
    },
  });

  res.status(200).send({ id: result.id, success: true });
};

/**
 * Updates a playlist with the given ID in the database.
 *
 * @param {Request} req - The HTTP request object.
 * @param {Response} res - The HTTP response object.
 * @param {User} user - The authenticated user object.
 * @return {Promise<void>} - A promise that resolves when the playlist is successfully updated.
 */
const updatePlaylist: AuthRouteHandlerMethod = async (req, res, user) => {
  const { id } = req.params as { id: string };
  const { name, shortDesc, poster } = req.body as { [key: string]: string };

  if (id.length === 0 || !validUUID(id)) {
    throw new Error("Invalid ID");
  }

  if (name.length === 0 || shortDesc.length === 0) {
    throw new Error("Name and Short Description is required");
  }

  const uploadedPoster = await cloudinary.uploader.upload(poster);

  const result = await prisma.playlist.update({
    data: {
      name,
      shortDesc,
      posterUrl: uploadedPoster.secure_url,
    },
    where: { id },
    select: {
      id: true,
    },
  });

  res.status(200).send({ id: result.id, success: true });
};

export {
  getPlaylistById,
  getPlaylists,
  likePlaylist,
  unlikePlaylist,
  uploadMusicToPlaylist,
  removeMusicFromPlaylist,
  createPlaylist,
  updatePlaylist,
};
