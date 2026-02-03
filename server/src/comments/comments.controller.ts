import { Controller, Delete, Param, UseGuards, Request, HttpCode } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { CommentsService } from './comments.service';
import { AuthGuard } from '@nestjs/passport';

@ApiTags('Comments')
@Controller('comments')
export class CommentsController {
  constructor(private readonly commentsService: CommentsService) {}

  @Delete(':commentId')
  @UseGuards(AuthGuard('jwt'))
  @ApiBearerAuth()
  @HttpCode(204)
  @ApiOperation({ summary: 'Delete a comment' })
  @ApiParam({ name: 'commentId', required: true })
  async deleteComment(@Request() req: any, @Param('commentId') commentId: string) {
    await this.commentsService.deleteComment(req.user.id, commentId);
  }
}

