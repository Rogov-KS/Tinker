/**
 * Seed script for demo data.
 *
 * Instructor usage:
 * 1. Ensure DATABASE_URL points to the target database.
 * 2. Run `npx prisma migrate deploy` (or `prisma migrate dev`) to apply schema.
 * 3. Execute `npx prisma db seed` to insert sample users, follows, posts and comments.
 */
import bcrypt from 'bcryptjs'
import { PrismaClient } from '@prisma/client'

const prisma = new PrismaClient()

async function main() {
  await prisma.like.deleteMany()
  await prisma.comment.deleteMany()
  await prisma.post.deleteMany()
  await prisma.follow.deleteMany()
  await prisma.user.deleteMany()

  const passwordHash = await bcrypt.hash('Password123!', 10)

  const [seneca, caesar, cicero, aurelius, plato] = await Promise.all([
    prisma.user.create({
      data: {
        username: 'seneca',
        password: passwordHash,
        bio: 'Римский философ-стоик и государственный деятель.',
        avatar:
          'https://upload.wikimedia.org/wikipedia/commons/5/5b/Seneca_statue.jpg',
      },
    }),
    prisma.user.create({
      data: {
        username: 'caesar',
        password: passwordHash,
        bio: 'Политик и полководец, автор «Записок о Галльской войне».',
        avatar:
          'https://upload.wikimedia.org/wikipedia/commons/5/5d/Gaius_Iulius_Caesar_%28Vatican_Museum%29.jpg',
      },
    }),
    prisma.user.create({
      data: {
        username: 'cicero',
        password: passwordHash,
        bio: 'Оратор и философ поздней Римской Республики.',
        avatar:
          'https://upload.wikimedia.org/wikipedia/commons/d/d0/M-T-Cicero.jpg',
      },
    }),
    prisma.user.create({
      data: {
        username: 'aurelius',
        password: passwordHash,
        bio: 'Император-философ, автор «Наедине с собой».',
        avatar:
          'https://upload.wikimedia.org/wikipedia/commons/9/9f/Marcus_Aurelius_Glyptothek_Munich.jpg',
      },
    }),
    prisma.user.create({
      data: {
        username: 'plato',
        password: passwordHash,
        bio: 'Древнегреческий философ, основатель Академии.',
        avatar:
          'https://upload.wikimedia.org/wikipedia/commons/4/4d/Plato-raphael.jpg',
      },
    }),
  ])

  await prisma.follow.createMany({
    data: [
      // seneca: 3 подписки (caesar, cicero, aurelius)
      { followerId: seneca.id, followingId: caesar.id },
      { followerId: seneca.id, followingId: cicero.id },
      { followerId: seneca.id, followingId: aurelius.id },
      // caesar: 2 подписки (seneca, cicero)
      { followerId: caesar.id, followingId: seneca.id },
      { followerId: caesar.id, followingId: cicero.id },
      // cicero: 2 подписки (seneca, aurelius)
      { followerId: cicero.id, followingId: seneca.id },
      { followerId: cicero.id, followingId: aurelius.id },
      // aurelius: 2 подписки (seneca, caesar)
      { followerId: aurelius.id, followingId: seneca.id },
      { followerId: aurelius.id, followingId: caesar.id },
      // plato: 0 подписок (не подписан ни на кого)
    ],
  })

  const posts = await prisma.$transaction([
    prisma.post.create({
      data: {
        userId: caesar.id,
        content:
          'Gallia est omnis divisa in partes tres. Каждый день учит нас быть стратегом.',
      },
    }),
    prisma.post.create({
      data: {
        userId: cicero.id,
        content:
          'Благо государства — высший закон. Не забывайте о добродетели в мелочах.',
      },
    }),
    prisma.post.create({
      data: {
        userId: aurelius.id,
        content:
          'Живи в согласии с природой. Лучшее время начать — сегодняшний день.',
      },
    }),
    prisma.post.create({
      data: {
        userId: seneca.id,
        content:
          'Всё, что происходит, может быть воспринято мудро. Дышите глубже.',
      },
    }),
    prisma.post.create({
      data: {
        userId: seneca.id,
        content:
          'Цель мудрости — соответствовать себе. Пишите посты только по делу.',
      },
    }),
  ])

  await prisma.comment.createMany({
    data: [
      {
        userId: seneca.id,
        postId: posts[0].id,
        content: 'Гай, стратегия впечатляет!',
      },
      {
        userId: aurelius.id,
        postId: posts[3].id,
        content: 'Спокойствие — союзник разума.',
      },
      {
        userId: cicero.id,
        postId: posts[4].id,
        content: 'Добродетель проявляется и в словах.',
      },
    ],
  })

  console.log('Seed data has been successfully inserted.')
}

main()
  .catch((error) => {
    console.error('Seed failed:', error)
    process.exit(1)
  })
  .finally(async () => {
    await prisma.$disconnect()
  })


