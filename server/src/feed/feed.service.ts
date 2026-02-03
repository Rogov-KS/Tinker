import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma.service';

@Injectable()
export class FeedService {
  constructor(private prisma: PrismaService) {}

  async getFeed(userId: string) {
    const following = await this.prisma.follow.findMany({
      where: { followerId: userId },
      select: { followingId: true },
    });
    
    const followingIds = following.map(f => f.followingId);

    const posts = await this.prisma.post.findMany({
      where: {
        userId: { in: followingIds },
      },
      include: {
        user: true,
        comments: {
           include: { user: true },
           orderBy: { createdAt: 'desc' }
        },
        likes: true,
      },
      orderBy: { createdAt: 'desc' },
    });

    return posts.map(post => this.mapPost(post, userId));
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
}

