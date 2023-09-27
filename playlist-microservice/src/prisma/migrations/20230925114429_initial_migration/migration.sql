-- CreateTable
CREATE TABLE "Playlist" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "likes" TEXT[],
    "musics" TEXT[],
    "posterUrl" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "shortDesc" TEXT NOT NULL,

    CONSTRAINT "Playlist_pkey" PRIMARY KEY ("id")
);
