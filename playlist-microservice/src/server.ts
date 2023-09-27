import Fastify from "fastify";
import dotenv from "dotenv";
import {
  getPlaylistById,
  getPlaylists,
  likePlaylist,
  removeMusicFromPlaylist,
  unlikePlaylist,
  uploadMusicToPlaylist,
  createPlaylist,
  updatePlaylist,
} from "./handlers";
import { auth, isOwner } from "./middlewares";
import fastifyMultipart from "@fastify/multipart";

const onFile = async (part: any) => {
  const buff = await part.toBuffer();
  const decoded = Buffer.from(buff.toString(), "base64").toString();
  part.value = decoded;
};

dotenv.config();

const fastify = Fastify({
  logger: true,
});

fastify.register(fastifyMultipart, { attachFieldsToBody: "keyValues", onFile });

fastify.setErrorHandler((err, req, res) => {
  res.status(500).send({
    success: false,
    error: err.message,
    stack: err.stack,
  });
});

fastify.get("/", (req, res) => {
  res.send({
    isActive: true,
  });
});

fastify.get("/getPlaylistById/:id", getPlaylistById);
fastify.get("/getPlaylists", getPlaylists);
fastify.post("/likePlaylist/:id", auth(likePlaylist));
fastify.post("/unlikePlaylist/:id", auth(unlikePlaylist));
fastify.post(
  "/uploadMusicToPlaylist/:id",
  auth(isOwner(uploadMusicToPlaylist))
);
fastify.put(
  "/removeMusicFromPlaylist/:id",
  auth(isOwner(removeMusicFromPlaylist))
);
fastify.post("/createPlaylist", auth(createPlaylist));
fastify.put("/updatePlaylist/:id", auth(isOwner(updatePlaylist)));

const port = process.env.PORT || 3000;

fastify.listen(
  {
    port: +port,
  },
  (err) => err && fastify.log.error(err)
);
