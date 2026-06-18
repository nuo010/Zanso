<template>
  <div class="landing-page">
    <nav>
      <div class="container">
        <a href="#" class="logo" @click.prevent="scrollToTop"><span>Z</span>anso</a>
        <div class="nav-links">
          <a href="#" @click.prevent="scrollToSection('features')">服务能力</a>
          <a href="#" @click.prevent="scrollToSection('advantages')">商业价值</a>
          <a href="#" @click.prevent="scrollToSection('usecases')">适合谁用</a>
          <a :href="appUrl" class="btn btn-primary nav-action" @click.prevent="goApp">立即开通</a>
        </div>
      </div>
    </nav>

    <section class="hero">
      <div class="container">
        <h1>把你的作品、产品和案例<br><span>变成更好成交的展示页</span></h1>
        <p>Zanso 帮你快速搭建专属数字展册，把图片、视频、项目资料整理成可分享、可交付、可持续复用的线上展示服务。客户看得清楚，沟通更高效，成交自然更利索。</p>
        <div class="hero-actions">
          <a :href="appUrl" class="btn btn-primary btn-hero" @click.prevent="goApp">马上创建展册</a>
          <a href="#" class="btn btn-outline btn-hero" @click.prevent="scrollToSection('advantages')">看看能赚在哪</a>
        </div>
      </div>
    </section>

    <section id="features">
      <div class="container text-center">
        <span class="section-label">服务能力</span>
        <h2 class="section-title">把展示、交付、获客放进一套系统</h2>
        <p class="section-subtitle">不是简单存文件，而是帮你把内容包装成客户愿意看的商业展示。</p>
      </div>
      <div class="container">
        <div class="feature-grid">
          <div class="feature-card">
            <div class="feature-icon feature-blue">📖</div>
            <h3>专属数字展册</h3>
            <p>为客户、项目、产品系列创建独立展册，封面、描述、分类清清楚楚。发出去不是一堆文件，而是一份像样的线上方案。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon feature-purple">🗂️</div>
            <h3>客户视角分类</h3>
            <p>按行业、项目阶段、产品型号或服务套餐整理内容，让客户不用翻聊天记录找素材，第一眼就知道你能提供什么。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon feature-green">🖼️</div>
            <h3>图片视频统一呈现</h3>
            <p>作品图、案例视频、产品细节、交付资料统一管理，批量上传、拖拽排序。展示效果越专业，客户信任来得越快。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon feature-orange">🔗</div>
            <h3>一键生成客户链接</h3>
            <p>给客户一个链接就能看完整内容，免登录、少打扰、方便转发。适合报价前展示、成交后交付、老客户复购触达。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon feature-red">📊</div>
            <h3>经营数据看板</h3>
            <p>资源数量、空间占用、分享内容集中查看，团队知道哪些内容在沉淀，老板也能看到服务资产是不是越攒越值钱。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon feature-cyan">👥</div>
            <h3>团队账号与权限</h3>
            <p>管理员统一维护用户和权限，销售、设计、运营各管各的内容，服务标准统一，客户体验也不再看个人发挥。</p>
          </div>
        </div>
      </div>
    </section>

    <section id="advantages" class="alt">
      <div class="container text-center">
        <span class="section-label">商业价值</span>
        <h2 class="section-title">让内容资产真正帮你赚钱</h2>
        <p class="section-subtitle">少一点反复解释，多一点专业呈现，把每一次分享都变成成交机会。</p>
      </div>
      <div class="container">
        <div class="why-grid">
          <div class="why-item">
            <span class="num">1</span>
            <h4>提升成交前的信任感</h4>
            <p>客户点开就是完整、整洁、有结构的内容展示，比微信群甩图、网盘丢文件专业太多，报价也更有底气。</p>
          </div>
          <div class="why-item">
            <span class="num">2</span>
            <h4>减少销售和交付沟通成本</h4>
            <p>资料统一入口，客户自己看、团队自己补，少来回翻文件、少重复介绍，时间省下来就能接更多单。</p>
          </div>
          <div class="why-item">
            <span class="num">3</span>
            <h4>沉淀可复用的服务资产</h4>
            <p>案例、素材、产品资料不是做完就散落各处，而是持续沉淀成可复用展册，新客户来了直接组合交付。</p>
          </div>
          <div class="why-item">
            <span class="num">4</span>
            <h4>适合做会员服务和客户交付</h4>
            <p>你可以把展册作为付费服务的一部分，按客户、套餐或项目交付专属内容，让服务看得见、摸得着。</p>
          </div>
          <div class="why-item">
            <span class="num">5</span>
            <h4>权限可控，分享更放心</h4>
            <p>展册和分类可控制可见性，分享链接可按需使用。公开展示、内部评审、客户交付都能分开管。</p>
          </div>
          <div class="why-item">
            <span class="num">6</span>
            <h4>稳定技术底座，能长期经营</h4>
            <p>基于 Go、Vue、MySQL、Redis 和对象存储构建，适合持续扩展，不是临时页面糊一下就完事。</p>
          </div>
        </div>
      </div>
    </section>

    <section id="usecases">
      <div class="container text-center">
        <span class="section-label">适合谁用</span>
        <h2 class="section-title">越依赖内容展示赚钱，越需要 Zanso</h2>
        <p class="section-subtitle">服务商、工作室、销售团队、教育机构，都可以把内容展示变成标准化交付能力。</p>
      </div>
      <div class="container">
        <div class="usecase-grid">
          <div class="usecase-card">
            <span class="emoji">🎨</span>
            <h4>设计与创意服务</h4>
            <p>把作品集、方案稿、交付案例整理成客户专属链接，提案更顺，报价更像回事。</p>
          </div>
          <div class="usecase-card">
            <span class="emoji">📸</span>
            <h4>摄影与影像交付</h4>
            <p>按客户和拍摄主题交付图片视频，选片、复购、转介绍都更方便。</p>
          </div>
          <div class="usecase-card">
            <span class="emoji">🏷️</span>
            <h4>产品目录与招商</h4>
            <p>产品图、视频、卖点资料统一展示，发给代理、客户、渠道，少讲废话多看内容。</p>
          </div>
          <div class="usecase-card">
            <span class="emoji">📚</span>
            <h4>课程与资料服务</h4>
            <p>把课程资料、案例库、会员素材按栏目交付，付费用户拿到的是完整服务，不是零散文件。</p>
          </div>
          <div class="usecase-card">
            <span class="emoji">🏗️</span>
            <h4>项目型服务交付</h4>
            <p>装修、工程、咨询、营销项目都能按阶段归档交付，让客户随时看到进展和成果。</p>
          </div>
        </div>
      </div>
    </section>

    <section id="cta" class="cta">
      <div class="container">
        <h2>现在就把内容变成可成交的服务资产</h2>
        <p>先创建一个客户展册，把最能打的案例、产品和交付资料放进去，下一次沟通直接发链接。</p>
        <a :href="appUrl" class="btn btn-cta btn-hero" @click.prevent="goApp">立即开通使用</a>
      </div>
    </section>

    <footer>
      <div class="container">
        <a href="#" class="logo" @click.prevent="scrollToTop"><span>Z</span>anso</a>
        <span>&copy; 2026 Zanso. All rights reserved.</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const appUrl = computed(() => {
  const baseUrl = `${window.location.origin}${window.location.pathname}`;
  return `${baseUrl.replace(/#.*$/, '')}#/login`;
});

function goApp() {
  router.push('/login');
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' });
}

function scrollToSection(id: string) {
  document.getElementById(id)?.scrollIntoView({
    behavior: 'smooth',
    block: 'start',
  });
}
</script>

<style scoped lang="scss">
.landing-page {
  min-height: 100vh;
  color: #1a1f36;
  background: #fafcff;
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
  --brand: #2463eb;
  --brand-dark: #1a4fc9;
  --accent: #7c5cfc;
  --text: #1a1f36;
  --text-secondary: #5e6f8d;
  --border: #e2e8f0;
  --card-bg: #ffffff;
  --section-bg: #f6f9ff;
}

.container {
  max-width: 1180px;
  margin: 0 auto;
  padding: 0 24px;
}

nav {
  position: sticky;
  top: 0;
  z-index: 100;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.85);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(16px);
}

nav .container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text);
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.5px;
  text-decoration: none;
}

.logo span {
  color: var(--brand);
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 32px;
}

.nav-links a {
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  transition: color 0.2s;
}

.nav-links a:hover {
  color: var(--brand);
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 22px;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-action {
  padding: 8px 18px;
  font-size: 13px;
}

.btn-primary {
  color: #fff !important;
  background: var(--brand);
}

.btn-primary:hover {
  color: #fff !important;
  background: var(--brand-dark);
  box-shadow: 0 8px 24px rgba(36, 99, 235, 0.25);
  transform: translateY(-1px);
}

.btn-outline {
  color: var(--text);
  background: transparent;
  border: 1.5px solid var(--border);
}

.btn-outline:hover {
  color: var(--brand);
  border-color: var(--brand);
}

.btn-hero {
  padding: 14px 32px;
  border-radius: 12px;
  font-size: 16px;
}

.hero {
  padding: 100px 0 80px;
  text-align: center;
  background: linear-gradient(180deg, #f6f9ff 0%, #ffffff 60%);
}

.hero h1 {
  max-width: 840px;
  margin: 0 auto 24px;
  font-size: clamp(40px, 6vw, 68px);
  font-weight: 800;
  line-height: 1.12;
  letter-spacing: -1.5px;
}

.hero h1 span {
  background: linear-gradient(135deg, var(--brand) 0%, var(--accent) 100%);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero p {
  max-width: 640px;
  margin: 0 auto 40px;
  color: var(--text-secondary);
  font-size: 18px;
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: center;
}

section {
  padding: 88px 0;
}

section.alt {
  background: var(--section-bg);
}

.section-label {
  display: inline-block;
  margin-bottom: 12px;
  color: var(--brand);
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 1.2px;
  text-transform: uppercase;
}

.section-title {
  margin-bottom: 12px;
  font-size: clamp(28px, 3.2vw, 40px);
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.8px;
}

.section-subtitle {
  max-width: 600px;
  margin: 0 auto 52px;
  color: var(--text-secondary);
  font-size: 17px;
  line-height: 1.7;
}

.text-center {
  text-align: center;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 28px;
}

.feature-card {
  padding: 36px 30px 32px;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  transition: transform 0.25s, box-shadow 0.25s;
}

.feature-card:hover {
  box-shadow: 0 24px 48px -12px rgba(0, 0, 0, 0.08);
  transform: translateY(-4px);
}

.feature-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  margin-bottom: 20px;
  border-radius: 12px;
  font-size: 22px;
}

.feature-blue {
  color: #2463eb;
  background: #eef4ff;
}

.feature-purple {
  color: #7c5cfc;
  background: #f3eeff;
}

.feature-green {
  color: #0d9488;
  background: #ecfdf5;
}

.feature-orange {
  color: #d97706;
  background: #fff7ed;
}

.feature-red {
  color: #dc2626;
  background: #fef2f2;
}

.feature-cyan {
  color: #0284c7;
  background: #f0f9ff;
}

.feature-card h3 {
  margin-bottom: 10px;
  font-size: 18px;
  font-weight: 600;
}

.feature-card p {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.65;
}

.why-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 24px;
}

.why-item {
  position: relative;
  padding: 28px 24px 28px 52px;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 14px;
}

.why-item .num {
  position: absolute;
  top: 28px;
  left: 20px;
  color: var(--brand);
  font-size: 20px;
  font-weight: 700;
}

.why-item h4 {
  margin-bottom: 6px;
  font-size: 15px;
  font-weight: 600;
}

.why-item p {
  color: var(--text-secondary);
  font-size: 13.5px;
  line-height: 1.6;
}

.usecase-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.usecase-card {
  padding: 28px 22px;
  text-align: center;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 14px;
}

.usecase-card .emoji {
  display: block;
  margin-bottom: 12px;
  font-size: 36px;
}

.usecase-card h4 {
  margin-bottom: 6px;
  font-size: 15px;
  font-weight: 600;
}

.usecase-card p {
  color: var(--text-secondary);
  font-size: 13px;
}

.cta {
  padding: 80px 0;
  text-align: center;
  background: linear-gradient(135deg, var(--brand) 0%, var(--accent) 100%);
}

.cta h2 {
  margin-bottom: 12px;
  color: #fff;
  font-size: clamp(28px, 3.2vw, 40px);
  font-weight: 700;
}

.cta p {
  margin-bottom: 32px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 17px;
}

.btn-cta {
  color: var(--brand);
  background: #fff;
}

.btn-cta:hover {
  color: var(--brand);
  background: #f4f8ff;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  transform: translateY(-1px);
}

footer {
  padding: 40px 0;
  border-top: 1px solid var(--border);
}

footer .container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: var(--text-secondary);
  font-size: 13px;
}

footer .logo {
  font-size: 16px;
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .hero {
    padding: 60px 0 48px;
  }

  section {
    padding: 56px 0;
  }

  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .hero-actions .btn {
    justify-content: center;
    width: 100%;
    max-width: 280px;
  }

  .feature-grid {
    gap: 20px;
  }

  .why-item {
    padding-left: 44px;
  }

  .why-item .num {
    left: 16px;
  }

  footer .container {
    flex-direction: column;
    gap: 12px;
    text-align: center;
  }
}
</style>
