import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from './ui/dialog';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormMessage,
} from './ui/form';
import { createUser } from '@/api';
import { useContext, useState } from 'react';
import { AuthContext } from '@/pages/AuthContext';
import { Loader2Icon } from 'lucide-react';

const FormSchema = z.object({
  name: z.string().min(2, {
    message: 'Username must be at least 2 characters.',
  }),
});

type Props = {
  onSuccess?: () => void;
  defaultOpen?: boolean;
  open?: boolean;
  isTriggerVisible?: boolean;
};

function AuthDialog({ onSuccess, open, defaultOpen, isTriggerVisible }: Props) {
  const { onSuccessfullUserCreate } = useContext(AuthContext);
  const [loading, setLoading] = useState(false);
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: '',
    },
  });

  async function onSubmit(data: z.infer<typeof FormSchema>) {
    setLoading(true);
    const userResponse = await createUser(data);
    setLoading(false);
    if (!userResponse) {
      toast.error('Could not create user');
      return;
    }
    onSuccessfullUserCreate(userResponse, onSuccess);
  }
  return (
    <Dialog defaultOpen={defaultOpen} open={open}>
      {isTriggerVisible && <DialogTrigger>Login</DialogTrigger>}
      <DialogContent
        className="sm:max-w-md"
        onEscapeKeyDown={(e) => e.preventDefault()}
        showCloseButton={false}
        onFocusOutside={(e) => e.preventDefault()}
        onPointerDownOutside={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle>Name yourself stranger</DialogTitle>
        </DialogHeader>
        <div className="flex items-center gap-2">
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="w-full space-y-6"
            >
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input placeholder="Your name" {...field} />
                    </FormControl>
                    <FormDescription>
                      This will be your public display name.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooter>
                <Button
                  type="submit"
                  className="align-self-end"
                  disabled={loading}
                >
                  {loading && <Loader2Icon className="animate-spin" />}
                  Continue
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default AuthDialog;
