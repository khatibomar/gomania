-- migrate:up
-- Seed categories
INSERT INTO
    categories (id, name, description, slug, is_active)
VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'Technology',
        'Tech news, programming, and digital innovation',
        'technology',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'Business',
        'Entrepreneurship, finance, and business strategies',
        'business',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440003',
        'Education',
        'Learning, teaching, and educational content',
        'education',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440004',
        'Entertainment',
        'Comedy, music, and entertainment shows',
        'entertainment',
        true
    ),
    (
        '550e8400-e29b-41d4-a716-446655440005',
        'News & Politics',
        'Current events and political discussions',
        'news-politics',
        true
    );

-- Seed users (CMS editors)
INSERT INTO
    users (
        id,
        email,
        password_hash,
        first_name,
        last_name,
        role,
        is_active
    )
VALUES
    (
        '660e8400-e29b-41d4-a716-446655440001',
        'admin@gomania.com',
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeGz4E7KmXuEQ8mSK',
        'Admin',
        'User',
        'admin',
        true
    ),
    (
        '660e8400-e29b-41d4-a716-446655440002',
        'editor@gomania.com',
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeGz4E7KmXuEQ8mSK',
        'Content',
        'Editor',
        'editor',
        true
    );

-- Seed programs
INSERT INTO
    programs (
        id,
        title,
        description,
        summary,
        category_id,
        language,
        country,
        author,
        publisher,
        artwork_url,
        website_url,
        is_explicit,
        status,
        total_episodes,
        average_duration,
        rating,
        total_ratings,
        source,
        created_by,
        published_at
    )
VALUES
    (
        '770e8400-e29b-41d4-a716-446655440001',
        'تقنية بودكاست',
        'برنامج أسبوعي يناقش أحدث التطورات في عالم التكنولوجيا والبرمجة، يستضيف خبراء ومطورين من المنطقة العربية لمناقشة التقنيات الحديثة والابتكارات.',
        'برنامج تقني أسبوعي يناقش أحدث التطورات التكنولوجية',
        '550e8400-e29b-41d4-a716-446655440001',
        'ar',
        'SA',
        'أحمد محمد',
        'شبكة تقنية',
        'https://example.com/artwork/tech-podcast.jpg',
        'https://techpodcast.sa',
        false,
        'active',
        25,
        1800,
        4.5,
        150,
        'local',
        '660e8400-e29b-41d4-a716-446655440001',
        '2024-01-15 10:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440002',
        'ريادة الأعمال العربية',
        'برنامج يسلط الضوء على قصص رواد الأعمال العرب ونجاحاتهم، مع تقديم نصائح عملية للمؤسسين الجدد ومناقشة التحديات والفرص في السوق العربي.',
        'قصص وتجارب رواد الأعمال في المنطقة العربية',
        '550e8400-e29b-41d4-a716-446655440002',
        'ar',
        'AE',
        'فاطمة السالم',
        'مؤسسة ريادة',
        'https://example.com/artwork/business-podcast.jpg',
        'https://entrepreneurship.ae',
        false,
        'active',
        18,
        2100,
        4.7,
        89,
        'local',
        '660e8400-e29b-41d4-a716-446655440002',
        '2024-02-01 14:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440003',
        'علوم المستقبل',
        'برنامج علمي يستكشف أحدث الاكتشافات العلمية والتقنيات المستقبلية، يناقش تأثيرها على حياتنا اليومية ومستقبل البشرية.',
        'اكتشافات علمية وتقنيات المستقبل',
        '550e8400-e29b-41d4-a716-446655440003',
        'ar',
        'EG',
        'د. محمد العلي',
        'أكاديمية العلوم',
        'https://example.com/artwork/science-podcast.jpg',
        'https://futurescience.eg',
        false,
        'active',
        32,
        1650,
        4.3,
        201,
        'local',
        '660e8400-e29b-41d4-a716-446655440001',
        '2024-01-20 16:30:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440004',
        'كوميديا الشارع',
        'برنامج كوميدي خفيف يناقش مواقف الحياة اليومية بطريقة فكاهية، مع ضيوف من عالم الكوميديا والفن.',
        'برنامج كوميدي يناقش الحياة اليومية',
        '550e8400-e29b-41d4-a716-446655440004',
        'ar',
        'JO',
        'سامر الكوميدي',
        'استوديو الضحك',
        'https://example.com/artwork/comedy-podcast.jpg',
        'https://streetcomedy.jo',
        false,
        'active',
        45,
        1200,
        4.1,
        312,
        'local',
        '660e8400-e29b-41d4-a716-446655440002',
        '2024-01-10 12:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440005',
        'أخبار التقنية اليومية',
        'نشرة إخبارية يومية تغطي أهم أخبار التكنولوجيا والذكاء الاصطناعي على مستوى العالم والمنطقة العربية.',
        'نشرة أخبار تقنية يومية',
        '550e8400-e29b-41d4-a716-446655440005',
        'ar',
        'SA',
        'نورا التقني',
        'شبكة الأخبار التقنية',
        'https://example.com/artwork/tech-news.jpg',
        'https://technews.sa',
        false,
        'active',
        156,
        900,
        4.6,
        543,
        'local',
        '660e8400-e29b-41d4-a716-446655440001',
        '2024-01-01 08:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440006',
        'تعلم البرمجة',
        'برنامج تعليمي يهدف إلى تعليم البرمجة للمبتدئين والمتقدمين، يغطي لغات البرمجة المختلفة والمفاهيم الأساسية.',
        'دروس برمجة للمبتدئين والمحترفين',
        '550e8400-e29b-41d4-a716-446655440003',
        'ar',
        'LB',
        'كريم المطور',
        'أكاديمية الكود',
        'https://example.com/artwork/programming.jpg',
        'https://learncode.lb',
        false,
        'active',
        28,
        2400,
        4.8,
        267,
        'local',
        '660e8400-e29b-41d4-a716-446655440002',
        '2024-02-15 10:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440007',
        'صوت الشباب',
        'برنامج يناقش قضايا الشباب العربي، التحديات والطموحات، مع استضافة شباب من مختلف البلدان العربية.',
        'برنامج يناقش قضايا وتطلعات الشباب العربي',
        '550e8400-e29b-41d4-a716-446655440005',
        'ar',
        'MA',
        'ليلى الشابة',
        'صوت الجيل',
        'https://example.com/artwork/youth-voice.jpg',
        'https://youthvoice.ma',
        false,
        'active',
        22,
        1950,
        4.2,
        134,
        'local',
        '660e8400-e29b-41d4-a716-446655440001',
        '2024-03-01 15:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440008',
        'مستثمر ذكي',
        'برنامج يناقش استراتيجيات الاستثمار والتخطيط المالي، مع خبراء في الأسواق المالية والاستثمار.',
        'نصائح استثمارية وتخطيط مالي',
        '550e8400-e29b-41d4-a716-446655440002',
        'ar',
        'KW',
        'عبدالله المالي',
        'مجموعة الاستثمار الذكي',
        'https://example.com/artwork/smart-investor.jpg',
        'https://smartinvestor.kw',
        false,
        'active',
        16,
        2700,
        4.4,
        98,
        'local',
        '660e8400-e29b-41d4-a716-446655440002',
        '2024-02-20 13:30:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440009',
        'تاريخ وحضارة',
        'رحلة في التاريخ العربي والإسلامي، استكشاف للحضارات والشخصيات التاريخية المؤثرة.',
        'استكشاف التاريخ والحضارة العربية الإسلامية',
        '550e8400-e29b-41d4-a716-446655440003',
        'ar',
        'SY',
        'د. أحمد المؤرخ',
        'معهد التاريخ',
        'https://example.com/artwork/history.jpg',
        'https://history.sy',
        false,
        'active',
        38,
        2200,
        4.6,
        189,
        'local',
        '660e8400-e29b-41d4-a716-446655440001',
        '2024-01-25 17:00:00+00'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440010',
        'صحة ولياقة',
        'برنامج يركز على الصحة العامة واللياقة البدنية، مع أطباء ومدربين لياقة لتقديم نصائح صحية عملية.',
        'نصائح صحية ولياقة بدنية',
        '550e8400-e29b-41d4-a716-446655440003',
        'ar',
        'QA',
        'د. سارة الصحة',
        'مركز الصحة الشاملة',
        'https://example.com/artwork/health.jpg',
        'https://healthfitness.qa',
        false,
        'active',
        31,
        1800,
        4.5,
        223,
        'local',
        '660e8400-e29b-41d4-a716-446655440002',
        '2024-03-10 11:00:00+00'
    );

-- Seed some episodes for the first few programs
INSERT INTO
    episodes (
        id,
        program_id,
        title,
        description,
        summary,
        episode_number,
        season_number,
        duration,
        audio_url,
        file_size,
        mime_type,
        is_explicit,
        status,
        published_at
    )
VALUES
    (
        '880e8400-e29b-41d4-a716-446655440001',
        '770e8400-e29b-41d4-a716-446655440001',
        'مقدمة في الذكاء الاصطناعي',
        'في هذه الحلقة نناقش أساسيات الذكاء الاصطناعي وتطبيقاته في حياتنا اليومية، مع استعراض أحدث التطورات في هذا المجال.',
        'مقدمة شاملة عن الذكاء الاصطناعي',
        1,
        1,
        1845,
        'https://example.com/audio/tech-podcast-ep1.mp3',
        45231678,
        'audio/mpeg',
        false,
        'published',
        '2024-01-15 10:30:00+00'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440002',
        '770e8400-e29b-41d4-a716-446655440001',
        'مستقبل البرمجة',
        'نتحدث عن مستقبل البرمجة واللغات الجديدة، وكيف ستؤثر التقنيات الناشئة على طريقة كتابة الكود.',
        'مناقشة مستقبل البرمجة واللغات الجديدة',
        2,
        1,
        1920,
        'https://example.com/audio/tech-podcast-ep2.mp3',
        47856234,
        'audio/mpeg',
        false,
        'published',
        '2024-01-22 10:30:00+00'
    ),
    (
        '880e8400-e29b-41d4-a716-446655440003',
        '770e8400-e29b-41d4-a716-446655440002',
        'من فكرة إلى شركة ناجحة',
        'قصة ملهمة لرائد أعمال عربي وكيف حول فكرته البسيطة إلى شركة تقنية ناجحة تخدم ملايين المستخدمين.',
        'قصة نجاح ملهمة لرائد أعمال عربي',
        1,
        1,
        2156,
        'https://example.com/audio/business-podcast-ep1.mp3',
        53782945,
        'audio/mpeg',
        false,
        'published',
        '2024-02-01 14:30:00+00'
    );

-- Seed some tags
INSERT INTO
    tags (id, name)
VALUES
    ('990e8400-e29b-41d4-a716-446655440001', 'تقنية'),
    (
        '990e8400-e29b-41d4-a716-446655440002',
        'ذكاء اصطناعي'
    ),
    (
        '990e8400-e29b-41d4-a716-446655440003',
        'ريادة أعمال'
    ),
    ('990e8400-e29b-41d4-a716-446655440004', 'برمجة'),
    ('990e8400-e29b-41d4-a716-446655440005', 'تعليم'),
    ('990e8400-e29b-41d4-a716-446655440006', 'صحة'),
    ('990e8400-e29b-41d4-a716-446655440007', 'تاريخ'),
    ('990e8400-e29b-41d4-a716-446655440008', 'استثمار'),
    ('990e8400-e29b-41d4-a716-446655440009', 'كوميديا'),
    ('990e8400-e29b-41d4-a716-446655440010', 'أخبار');

-- migrate:down
DELETE FROM episodes
WHERE
    program_id IN (
        SELECT
            id
        FROM
            programs
        WHERE
            source = 'local'
            AND id LIKE '770e8400-e29b-41d4-a716-44665544%'
    );

DELETE FROM programs
WHERE
    source = 'local'
    AND id LIKE '770e8400-e29b-41d4-a716-44665544%';

DELETE FROM tags
WHERE
    id LIKE '990e8400-e29b-41d4-a716-44665544%';

DELETE FROM users
WHERE
    id LIKE '660e8400-e29b-41d4-a716-44665544%';

DELETE FROM categories
WHERE
    id LIKE '550e8400-e29b-41d4-a716-44665544%';
