import { Injectable, BadRequestException } from '@nestjs/common';
import { PrismaService } from '../prisma.service';

@Injectable()
export class SearchService {
  constructor(private prisma: PrismaService) {}

  async search(query: string, type: string) {
    if (type === 'users') {
      const users = await this.prisma.user.findMany({
        where: {
          username: { contains: query, mode: 'insensitive' },
        },
        include: { following: true },
      });
      return users.map(user => ({
        id: user.id,
        username: user.username,
        avatar: user.avatar,
        bio: user.bio,
        following: user.following?.map((f: any) => f.followingId) || [],
      }));
    } else if (type === 'posts') {
      const posts = await this.prisma.post.findMany({
        where: {
          content: { contains: query, mode: 'insensitive' },
        },
        include: {
          user: true,
          comments: { include: { user: true } },
          likes: true,
        },
        orderBy: { createdAt: 'desc' },
      });
      return posts.map(post => this.mapPost(post));
    } else {
        throw new BadRequestException('Invalid type');
    }
  }

  private mapPost(post: any) {
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
      likedByCurrentUser: false,
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
}

