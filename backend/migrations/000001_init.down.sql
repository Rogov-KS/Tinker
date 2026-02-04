-- DropForeignKey
ALTER TABLE "Follow" DROP CONSTRAINT IF EXISTS "Follow_followingId_fkey";
ALTER TABLE "Follow" DROP CONSTRAINT IF EXISTS "Follow_followerId_fkey";
ALTER TABLE "Like" DROP CONSTRAINT IF EXISTS "Like_postId_fkey";
ALTER TABLE "Like" DROP CONSTRAINT IF EXISTS "Like_userId_fkey";
ALTER TABLE "Comment" DROP CONSTRAINT IF EXISTS "Comment_postId_fkey";
ALTER TABLE "Comment" DROP CONSTRAINT IF EXISTS "Comment_userId_fkey";
ALTER TABLE "Post" DROP CONSTRAINT IF EXISTS "Post_userId_fkey";

-- DropIndex
DROP INDEX IF EXISTS "Follow_followerId_followingId_key";
DROP INDEX IF EXISTS "Follow_followingId_idx";
DROP INDEX IF EXISTS "Follow_followerId_idx";
DROP INDEX IF EXISTS "Like_userId_postId_key";
DROP INDEX IF EXISTS "Like_userId_idx";
DROP INDEX IF EXISTS "Like_postId_idx";
DROP INDEX IF EXISTS "Comment_userId_idx";
DROP INDEX IF EXISTS "Comment_postId_idx";
DROP INDEX IF EXISTS "Post_createdAt_idx";
DROP INDEX IF EXISTS "Post_userId_idx";
DROP INDEX IF EXISTS "User_username_key";

-- DropTable
DROP TABLE IF EXISTS "Follow";
DROP TABLE IF EXISTS "Like";
DROP TABLE IF EXISTS "Comment";
DROP TABLE IF EXISTS "Post";
DROP TABLE IF EXISTS "User";

