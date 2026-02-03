import { Controller, Get, Post, Body, Patch, Param, UseGuards, Request, Delete, HttpCode, ForbiddenException } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { UsersService } from './users.service';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { CreatePostDto } from '../posts/dto/create-post.dto';
import { AuthGuard } from '@nestjs/passport';

@ApiTags('Users')
@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Get()
  @ApiOperation({ summary: 'Get all users' })
  getAllUsers() {
    return this.usersService.getAllUsers();
  }

  @Post()
  @ApiOperation({ summary: 'Create a new user' })
  createUser(@Body() createUserDto: CreateUserDto) {
    return this.usersService.createUser(createUserDto);
  }

  @Get(':userId')
  @ApiOperation({ summary: 'Get user by ID' })
  @ApiParam({ name: 'userId', required: true })
  getUserById(@Param('userId') userId: string) {
    return this.usersService.getUserById(userId);
  }

  @Patch(':userId')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Update user profile' })
  @ApiParam({ name: 'userId', required: true })
  updateUser(@Request() req: any, @Param('userId') userId: string, @Body() updateUserDto: UpdateUserDto) {
    // Only allow users to update their own profile
    if (req.user.id !== userId) {
      throw new ForbiddenException('You can only update your own profile');
    }
    return this.usersService.updateUser(userId, updateUserDto);
  }

  @Get(':userId/posts')
  @ApiOperation({ summary: "Get all posts by user" })
  @ApiParam({ name: 'userId', required: true })
  getUserPosts(@Param('userId') userId: string) {
    return this.usersService.getUserPosts(userId);
  }

  @Post(':userId/posts')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Create a new post' })
  @ApiParam({ name: 'userId', required: true })
  createPost(@Request() req: any, @Param('userId') userId: string, @Body() createPostDto: CreatePostDto) {
    // Only allow users to create posts as themselves
    if (req.user.id !== userId) {
      throw new ForbiddenException('You can only create posts as yourself');
    }
    return this.usersService.createPost(userId, createPostDto);
  }

  @Post(':userId/follow')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Follow a user' })
  @ApiParam({ name: 'userId', required: true })
  async followUser(@Request() req: any, @Param('userId') userId: string) {
    await this.usersService.followUser(req.user.id, userId);
    return { message: 'User followed' };
  }

  @Delete(':userId/follow')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @HttpCode(204)
  @ApiOperation({ summary: 'Unfollow a user' })
  @ApiParam({ name: 'userId', required: true })
  async unfollowUser(@Request() req: any, @Param('userId') userId: string) {
    await this.usersService.unfollowUser(req.user.id, userId);
  }

  @Get(':userId/following')
  @ApiOperation({ summary: 'Get list of users being followed' })
  @ApiParam({ name: 'userId', required: true })
  getFollowing(@Param('userId') userId: string) {
    return this.usersService.getFollowing(userId);
  }
}

