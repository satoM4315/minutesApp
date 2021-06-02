import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles({
  root: {
    minWidth: 400,
    maxWidth: 800,
    marginTop: 20,
    marginBottom: 20,
    marginLeft: 20,
    marginRight: 20,
  },
  bullet: {
    display: 'inline-block',
    margin: '0 2px',
    transform: 'scale(0.8)',
  },
  title: {
    fontSize: 18,
  },
  pos: {
    marginBottom: 12,
  },
});

export default function GetImportantSentences() {
  const [data, setData] = useState<string[]>([]);
  const classes = useStyles();

  useEffect(() => {
    // ルート /message に対して GETリクエストを送る
    // 帰ってきたものをjsonにしてuseStateに突っ込む
    fetch(`${window.location.pathname}/important_sentences`)
      .then((res) => res.json())
      .then(setData);
  }, []);

  return (
    // タグが複数できる場合は何らかのタグで全体を囲う
    <div>
      <Card className={classes.root} variant="outlined">
        <CardContent>
          <Typography variant="h5" component="h2" className={classes.title}>
            議事録から抽出した重要と思われる発言一覧
          </Typography>
          <Typography variant="body2" component="div">
            {data.map((item) => (
              <p key={item}>{item}</p>
            ))}
          </Typography>
        </CardContent>
      </Card>
    </div>
  );
}
