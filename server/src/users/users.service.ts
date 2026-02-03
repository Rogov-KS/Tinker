import { Injectable, NotFoundException, BadRequestException } from '@nestjs/common';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';
import { PrismaService } from '../prisma.service';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { CreatePostDto } from '../posts/dto/create-post.dto';
import * as bcrypt from 'bcryptjs';

@Injectable()
export class UsersService {
  constructor(
    private prisma: PrismaService,
    @InjectPinoLogger(UsersService.name)
    private readonly logger: PinoLogger,
  ) {}

  async getAllUsers() {
    const users = await this.prisma.user.findMany({
      include: {
        following: true,
      },
      orderBy: { createdAt: 'asc' }, // Раньше зарегистрировался = выше
    });
    return users.map(user => this.mapUser(user));
  }

  async createUser(createUserDto: CreateUserDto) {
    this.logger.info(`Creating new user: ${createUserDto.username}`);
    
    const existing = await this.prisma.user.findUnique({
      where: { username: createUserDto.username },
    });
    if (existing) {
      this.logger.warn(`Username already exists: ${createUserDto.username}`);
      throw new BadRequestException('Username already exists');
    }

    const hashedPassword = await bcrypt.hash(createUserDto.password, 10);
    const user = await this.prisma.user.create({
      data: {
        username: createUserDto.username,
        password: hashedPassword,
      },
      include: { following: true },
    });
    
    this.logger.info({ userId: user.id, username: user.username }, `User created successfully`);
    return this.mapUser(user);
  }

  async getUserById(userId: string) {
    const user = await this.prisma.user.findUnique({
      where: { id: userId },
      include: { following: true },
    });
    if (!user) throw new NotFoundException('User not found');
    return this.mapUser(user);
  }

  async updateUser(userId: string, updateUserDto: UpdateUserDto) {
    const user = await this.prisma.user.update({
      where: { id: userId },
      data: updateUserDto,
      include: { following: true },
    });
    return this.mapUser(user);
  }

  async followUser(followerId: string, followingId: string) {
    if (followerId === followingId) {
      throw new BadRequestException('Cannot follow yourself');
    }
    
    // Check if already following
    const existing = await this.prisma.follow.findUnique({
      where: {
        followerId_followingId: {
          followerId,
          followingId,
        },
      },
    });

    if (existing) {
       // Ideally just return success or throw, but let's be idempotent
       return; 
    }

    await this.prisma.follow.create({
      data: {
        followerId,
        followingId,
      },
    });
  }

  async unfollowUser(followerId: string, followingId: string) {
    try {
      await this.prisma.follow.delete({
        where: {
          followerId_followingId: {
            followerId,
            followingId,
          },
        },
      });
    } catch (e) {
      // Ignore if not found
    }
  }

  async getFollowing(userId: string) {
    const follows = await this.prisma.follow.findMany({
      where: { followerId: userId },
      include: { following: { include: { following: true } } },
    });
    return follows.map(f => this.mapUser(f.following));
  }

  async getUserPosts(userId: string) {
    const posts = await this.prisma.post.findMany({
      where: { userId },
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

  async createPost(userId: string, createPostDto: CreatePostDto) {
    this.logger.info({ userId }, `Creating post for user`);
    
    const post = await this.prisma.post.create({
      data: {
        userId,
        content: createPostDto.content,
      },
      include: {
        user: true,
        comments: { include: { user: true } },
        likes: true,
      },
    });
    
    this.logger.info({ postId: post.id, userId }, `Post created`);
    return this.mapPost(post, userId);
  }

  private mapUser(user: any) {
    return {
      id: user.id,
      username: user.username,
      avatar: user.avatar,
      bio: user.bio,
      following: user.following?.map((f: any) => f.followingId) || [],
    };
  }

  public mapPost(post: any, currentUserId?: string) {
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

