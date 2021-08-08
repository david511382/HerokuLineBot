import Head from 'next/head'
import Image from 'next/image'
import styles from './Home.module.css'

export default function Home() {
  return (
    <div className={styles.container}>
      <main className={styles.main}>
        <h1 className={styles.title}>
          大台中臭豆腐商行
        </h1>

        <div className={styles.grid}>
          <a className={styles.card}>
            <h2>臭豆腐</h2>
            <p>適合油炸的大小</p>
          </a>

          <a className={styles.card}>
            <h2>麻辣臭豆腐</h2>
            <p>厚實的臭豆腐，適合煮麻辣</p>
          </a>

          <a className={styles.card}>
            <h2>炭烤臭豆腐</h2>
            <p>適合炭烤</p>
          </a>
        </div>
        
        <p className={styles.description}>
          訂購電話: 0936329626
        </p>
        <p className={styles.description}>
          公司電話: 0423376183
        </p>

        <p className={styles.description}>
          Line 帳號 掃描 QR Code 加入好友
        </p>
        <Image
          src="/line.png"
          alt="大台中臭豆腐 Line QR code"
          width={200}
          height={200} 
        />
      </main>
      
      <footer className={styles.footer}>
        <a>
          大台中臭豆腐商行 公司電話: 0423376183 手機: 0936329626
        </a>
      </footer>
    </div>
  )
}
