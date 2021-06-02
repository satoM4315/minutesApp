import React, { useState } from 'react';
import PropTypes from 'prop-types';
import ReactDOM from 'react-dom';
import { createStyles, makeStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Alert from '@material-ui/lab/Alert';

// メッセージ追加のAPIへのURL
// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_LOGIN = '/register';

const useStyles = makeStyles((theme) =>
  createStyles({
    root: {
      '& > *': {
        margin: theme.spacing(1),
        width: '25ch',
      },
      width: '100%',
      '& > * + *': {
        marginTop: theme.spacing(2),
      },
    },
  })
);

export default function RegisterPostForm(props) {
  // テキストボックス内のメッセージ
  const [userId, setUserId] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [passwordForConfirm, setPasswordForConfirm] = useState<string>('');
  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const classes = useStyles();

  const handleSubmit = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true);
    // 登録ボタンを押された際に判定(リアルタイムでやりたい)
    if (password === passwordForConfirm) {
      try {
        // ページが更新されないようにする
        event.preventDefault();

        // Reactのハンドラはasyncにできる
        const res = await fetch(API_URL_LOGIN, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          // 相応しくないかも
          // same-originを使うべき？
          credentials: 'include',
          body: JSON.stringify({ userId, password }),
        });

        setUserId('');
        setPassword('');
        setPasswordForConfirm('');

        const obj = await res.json();
        if ('error' in obj) {
          // サーバーからエラーが返却された
          ReactDOM.render(
            <div className={classes.root}>
              <Alert
                variant="outlined"
                severity="error"
                onClose={() => {
                  ReactDOM.render(
                    <div />,
                    document.getElementById('serverMessage')
                  );
                }}
              >
                {obj.error}
              </Alert>
            </div>,
            document.getElementById('serverMessage')
          );
          throw new Error(
            `An error occurred on querying ${API_URL_LOGIN}, the response included error message: ${obj.error}`
          );
        }
        if (!('success' in obj)) {
          // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
          ReactDOM.render(
            <div className={classes.root}>
              <Alert
                variant="outlined"
                severity="error"
                onClose={() => {
                  ReactDOM.render(
                    <div />,
                    document.getElementById('serverMessage')
                  );
                }}
              >
                Error
              </Alert>
            </div>,
            document.getElementById('serverMessage')
          );
          throw new Error(
            `An response from ${API_URL_LOGIN} unexpectedly did not have 'success' member`
          );
        }
        if (obj.success !== true) {
          ReactDOM.render(
            <div className={classes.root}>
              <Alert
                variant="outlined"
                severity="error"
                onClose={() => {
                  ReactDOM.render(
                    <div />,
                    document.getElementById('serverMessage')
                  );
                }}
              >
                Error
              </Alert>
            </div>,
            document.getElementById('serverMessage')
          );
          throw new Error(
            `An response from ${API_URL_LOGIN} returned non true value as 'success' member`
          );
        }

        // 要求は成功
        ReactDOM.render(
          <div className={classes.root}>
            <Alert variant="outlined" severity="success">
              登録完了! 3秒後にログインページへ推移
            </Alert>
            <CircularProgress />
          </div>,
          document.getElementById('serverMessage')
        );
        // リスナ関数を呼ぶ
        props.onSubmitSuccessful();

        // 登録が成功したらログインページにリダイレクト
        setTimeout(() => {
          window.location.href = '/login';
        }, 3000);
      } finally {
        setWorking(false);
      }
    } else {
      try {
        ReactDOM.render(
          <div className={classes.root}>
            <Alert
              variant="outlined"
              severity="error"
              onClose={() => {
                ReactDOM.render(
                  <div />,
                  document.getElementById('serverMessage')
                );
              }}
            >
              間違ってんで
            </Alert>
          </div>,
          document.getElementById('serverMessage')
        );
      } finally {
        setWorking(false);
      }
    }
  };

  return (
    <>
      <form className={classes.root} noValidate autoComplete="off">
        <TextField
          id="standard-basic"
          label="ユーザーID"
          value={userId}
          type="textbox"
          onChange={(event) => setUserId(event.target.value)}
        />
        <p />
        <TextField
          id="standard-basic"
          label="パスワード"
          value={password}
          type="password"
          onChange={(event) => setPassword(event.target.value)}
        />
        <p />
        <TextField
          id="standard-basic"
          label="確認用パスワード"
          value={passwordForConfirm}
          type="password"
          onChange={(event) => setPasswordForConfirm(event.target.value)}
          // onChange={(event) => comparePassword(event.target.value)}
        />
        <p />
        <Button
          variant="contained"
          color="primary"
          disabled={working}
          onClick={handleSubmit}
        >
          登録
        </Button>
      </form>
    </>
  );
  // var comparePassword: (passwordForConfirm: string) => return string===password;
}

RegisterPostForm.propTypes = {
  onSubmitSuccessful: PropTypes.func,
};

RegisterPostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};
