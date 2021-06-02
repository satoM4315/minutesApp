import React, { useState } from 'react';
import Button from '@material-ui/core/Button';
import PropTypes from 'prop-types';
import ReactDOM from 'react-dom';
import Alert from '@material-ui/lab/Alert';
import KeyboardVoiceIcon from '@material-ui/icons/KeyboardVoice';
import { createStyles, makeStyles } from '@material-ui/core/styles';

const apiUrlAddMessage = `${window.location.pathname}/add_message`;

const useStyles = makeStyles((theme) =>
  createStyles({
    button: {
      margin: theme.spacing(1),
    },
  })
);

export default function AudioMessagePostForm(props) {
  const [working, setWorking] = useState<boolean>(false);
  const [action, setAction] = useState<String>('Start');
  const classes = useStyles();

  if ('SpeechRecognition' in window) {
    (window as any).SpeechRecognition = (window as any).SpeechRecognition;
  } else if ('webkitSpeechRecognition' in window) {
    (window as any).SpeechRecognition = (window as any).webkitSpeechRecognition;
  } else {
    return (
      <Alert
        variant="outlined"
        severity="error"
        onClose={() => {
          ReactDOM.render(<div />, document.getElementById('serverMessage'));
        }}
      >
        このブラウザは音声認識に対応していません
      </Alert>
    );
  }

  const speech = new (window as any).SpeechRecognition();
  speech.lang = 'ja-JP';
  speech.onresult = async function AudioResult(e) {
    speech.stop();
    if (e.results[0].isFinal) {
      const audioText = e.results[0][0].transcript;
      const res = await fetch(apiUrlAddMessage, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ message: audioText }),
      });
      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        throw new Error(
          `An error occurred on querying ${apiUrlAddMessage}, the response included error message: ${obj.error}`
        );
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        throw new Error(
          `An response from ${apiUrlAddMessage} unexpectedly did not have 'success' member`
        );
      }
      if (obj.success !== true) {
        throw new Error(
          `An response from ${apiUrlAddMessage} returned non true value as 'success' member`
        );
      }
      props.onSubmitSuccessful();
    }
  };

  speech.onend = () => {
    speech.start();
  };

  const handleSubmit = (event: React.FormEvent) => {
    setAction('Recording Now');
    setWorking(true);
    try {
      // ページが更新されないようにする
      speech.start();
      event.preventDefault();
      props.onSubmitSuccessful();
    } finally {
      setWorking(false);
    }
  };
  return (
    <form>
      <Button
        disabled={working}
        variant="contained"
        color="primary"
        endIcon={<KeyboardVoiceIcon />}
        className={classes.button}
        onClick={handleSubmit}
      >
        {action}
      </Button>
      <Button
        disabled={working}
        variant="contained"
        color="secondary"
        className={classes.button}
        onClick={() => {
          setAction('Start');
          window.location.href = window.location.pathname;
          // props.onSubmitSuccessful();
        }}
      >
        Stop
      </Button>
    </form>
  );
}
AudioMessagePostForm.propTypes = {
  // 新しいメッセージの追加が正常に完了したら呼ばれる関数
  onSubmitSuccessful: PropTypes.func,
};

AudioMessagePostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};
