import { ForbiddenException, Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from '../prisma.service';

@Injectable()
export class PostsService {
  constructor(private prisma: PrismaService) {}

  async getAllPosts(currentUserId?: string) {
    const posts = await this.prisma.post.findMany({
      include: {
        user: true,
        comments: {
          include: { user: true },
          orderBy: { createdAt: 'desc' }
        },
        likes: true,
      },
      orderBy: { createdAt: 'desc' }, // Последний пост выше
    });

    return posts.map(post => this.mapPost(post, currentUserId));
  }

  private mapPost(post: any, currentUserId?: string) {
    return {
      id: post.id,
      userId: post.userId,
      author: {
        id: post.user.id,
        username: post.user.username,
        avatar: post.user.avatar,
      },
      content: post.content,
      likes: post.likes.length,
      likedByCurrentUser: currentUserId ? post.likes.some((l: any) => l.userId === currentUserId) : false,
      createdAt: post.createdAt,
      comments: post.comments.map((c: any) => ({
        id: c.id,
        userId: c.userId,
        content: c.content,
        createdAt: c.createdAt,
        author: {
          id: c.user.id,
          username: c.user.username,
          avatar: c.user.avatar,
        }
      })),
    };
  }

  async deletePost(userId: string, postId: string) {
    const post = await this.prisma.post.findUnique({ where: { id: postId } });
    if (!post) throw new NotFoundException('Post not found');
    if (post.userId !== userId) throw new ForbiddenException('You can only delete your own posts');
    
    await this.prisma.post.delete({ where: { id: postId } });
  }

  async likePost(userId: string, postId: string) {
    // Check if post exists
    const post = await this.prisma.post.findUnique({ where: { id: postId } });
    if (!post) throw new NotFoundException('Post not found');

    try {
      await this.prisma.like.create({
        data: { userId, postId },
      });
    } catch (e) {
      // Ignore unique constraint violation
    }
  }

  async unlikePost(userId: string, postId: string) {
    try {
      await this.prisma.like.delete({
        where: {
          userId_postId: { userId, postId },
        },
      });
    } catch (e) {
      // Ignore if not found
    }
  }

  async getComments(postId: string) {
    const comments = await this.prisma.comment.findMany({
      where: { postId },
      include: { user: true },
      orderBy: { createdAt: 'desc' },
    });
    return comments.map(c => ({
      id: c.id,
      userId: c.userId,
      content: c.content,
      createdAt: c.createdAt,
      author: {
        id: c.user.id,
        username: c.user.username,
        avatar: c.user.avatar,
      }
    }));
  }

  async createComment(userId: string, postId: string, content: string) {
    const post = await this.prisma.post.findUnique({ where: { id: postId } });
    if (!post) throw new NotFoundException('Post not found');

    const comment = await this.prisma.comment.create({
      data: {
        userId,
        postId,
        content,
      },
      include: { user: true },
    });
    return {
      id: comment.id,
      userId: comment.userId,
      content: comment.content,
      createdAt: comment.createdAt,
      author: {
        id: comment.user.id,
        username: comment.user.username,
        avatar: comment.user.avatar,
      }
    };
  }
}

