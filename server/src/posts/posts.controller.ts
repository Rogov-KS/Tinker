import { Controller, Post, Get, Delete, Body, Param, UseGuards, Request, HttpCode } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { PostsService } from './posts.service';
import { CreateCommentDto } from './dto/create-comment.dto';
import { AuthGuard } from '@nestjs/passport';

@ApiTags('Posts')
@Controller('posts')
export class PostsController {
  constructor(private readonly postsService: PostsService) {}

  @Get()
  @ApiOperation({ summary: 'Get all posts' })
  getAllPosts(@Request() req: any) {
    // Optional auth - if user is authenticated, pass userId for like status
    const currentUserId = req.user?.id;
    return this.postsService.getAllPosts(currentUserId);
  }

  @Delete(':postId')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @HttpCode(204)
  @ApiOperation({ summary: 'Delete a post' })
  @ApiParam({ name: 'postId', required: true })
  async deletePost(@Request() req: any, @Param('postId') postId: string) {
    await this.postsService.deletePost(req.user.id, postId);
  }

  @Post(':postId/likes')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Like a post' })
  @ApiParam({ name: 'postId', required: true })
  async likePost(@Request() req: any, @Param('postId') postId: string) {
    await this.postsService.likePost(req.user.id, postId);
    return { message: 'Post liked' };
  }

  @Delete(':postId/likes')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @HttpCode(204)
  @ApiOperation({ summary: 'Unlike a post' })
  @ApiParam({ name: 'postId', required: true })
  async unlikePost(@Request() req: any, @Param('postId') postId: string) {
    await this.postsService.unlikePost(req.user.id, postId);
  }

  @Get(':postId/comments')
  @ApiOperation({ summary: 'Get all comments for a post' })
  @ApiParam({ name: 'postId', required: true })
  getComments(@Param('postId') postId: string) {
    return this.postsService.getComments(postId);
  }

  @Post(':postId/comments')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Create a new comment' })
  @ApiParam({ name: 'postId', required: true })
  createComment(@Request() req: any, @Param('postId') postId: string, @Body() createCommentDto: CreateCommentDto) {
    return this.postsService.createComment(req.user.id, postId, createCommentDto.content);
  }
}

